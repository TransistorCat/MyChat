package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"google.golang.org/grpc"

	pb "varifyserver/pb"
)

type server struct {
	pb.UnimplementedVarifyServiceServer
}

// GetVarifyCode grpc响应获取验证码的服务
func (s *server) GetVarifyCode(ctx context.Context, req *pb.GetVarifyReq) (*pb.GetVarifyRsp, error) {
	fmt.Println("email is", req.GetEmail())
	var uniqueId string

	queryRes, err := getRedis(CodePrefix + req.GetEmail())
	if err != nil {
		log.Println("Redis query error:", err)
		return &pb.GetVarifyRsp{Email: req.GetEmail(), Error: RedisErr}, nil
	}

	if queryRes == "" {
		uniqueId = uuid.New().String()
		if len(uniqueId) > 4 {
			uniqueId = uniqueId[:4]
		}

		err := setRedisExpire(CodePrefix+req.GetEmail(), uniqueId, 15*time.Minute)
		if err != nil {
			log.Println("Redis set error:", err)
			return &pb.GetVarifyRsp{Email: req.GetEmail(), Error: RedisErr}, nil
		}
	} else {
		uniqueId = queryRes
	}

	fmt.Println("uniqueId is", uniqueId)
	textStr := fmt.Sprintf("您的验证码为%s，请十五分钟内完成注册", uniqueId)

	// 发送邮件
	mailOptions := &email.Email{
		From:    EmailUser,
		To:      []string{req.GetEmail()},
		Subject: "验证码",
		Text:    []byte(textStr),
	}

	sendRes, err := sendMail(mailOptions)
	if err != nil {
		log.Println("SendMail error:", err)
		return &pb.GetVarifyRsp{Email: req.GetEmail(), Error: Exception}, nil
	}

	fmt.Println("send res is", sendRes)
	return &pb.GetVarifyRsp{Email: req.GetEmail(), Error: Success}, nil
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterVarifyServiceServer(s, &server{})

	fmt.Println("varify server started")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
