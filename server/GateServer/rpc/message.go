// package rpc

// import (
// 	"GateServer/pb"
// 	"context"
// 	"flag"

// 	"google.golang.org/grpc"
// )

// var (
// 	addr = flag.String("addr", "localhost:50051", "the address to connect to")
// )

// type MessageManager struct {
// 	conn   *grpc.ClientConn
// 	client pb.VarifyServiceClient
// }

// // NewRPCManager 创建一个新的 RPCManager 实例并初始化 gRPC 客户端连接
// func NewMessageManager(addr string) (*MessageManager, error) {
// 	conn, err := grpc.NewClient(addr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client := pb.NewVarifyServiceClient(conn)
// 	return &MessageManager{
// 		conn:   conn,
// 		client: client,
// 	}, nil
// }

// func (m *MessageManager) GetVarifyCode(email string) {
// 	m.client.GetVarifyCode(context.Background(), &pb.GetVarifyReq{Email: email})
// }

// // Close 关闭 gRPC 客户端连接
// func (m *MessageManager) Close() error {
// 	if m.conn != nil {
// 		return m.conn.Close()
// 	}
// 	return nil
// }

package rpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "GateServer/pb" // Update with your actual proto package path
)

type RPConPool struct {
	poolSize    int
	host        string
	port        string
	connections chan pb.VarifyServiceClient
	mutex       sync.Mutex
}

func NewRPConPool(poolSize int, host, port string) *RPConPool {
	pool := &RPConPool{
		poolSize:    poolSize,
		host:        host,
		port:        port,
		connections: make(chan pb.VarifyServiceClient, poolSize),
	}

	for i := 0; i < poolSize; i++ {
		conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Failed to dial: %v\n", err)
			continue
		}
		client := pb.NewVarifyServiceClient(conn)
		pool.connections <- client
	}

	return pool
}

func (p *RPConPool) getConnection() (pb.VarifyServiceClient, error) {
	select {
	case conn := <-p.connections:
		return conn, nil
	default:
		conn, err := grpc.Dial(p.host+":"+p.port, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		return pb.NewVarifyServiceClient(conn), nil
	}
}

func (p *RPConPool) returnConnection(conn pb.VarifyServiceClient) {
	p.connections <- conn
}

type VerifyGrpcClient struct {
	pool *RPConPool
}

func NewVerifyGrpcClient(pool *RPConPool) *VerifyGrpcClient {
	return &VerifyGrpcClient{
		pool: pool,
	}
}

func (c *VerifyGrpcClient) GetVarifyCode(email string) (*pb.GetVarifyRsp, error) {
	conn, err := c.pool.getConnection()
	if err != nil {
		return nil, err
	}
	defer c.pool.returnConnection(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	request := &pb.GetVarifyReq{
		Email: email,
	}

	reply, err := conn.GetVarifyCode(ctx, request)
	if err != nil {
		return nil, status.Errorf(status.Code(err), "RPC failed: %v", err)
	}

	return reply, nil
}

// func main() {
// 	pool := NewRPConPool(5, "localhost", "50051") // Adjust host and port accordingly
// 	client := NewVerifyGrpcClient(pool)

// 	email := "1820737440@qq.com"
// 	response, err := client.GetVarifyCode(email)
// 	if err != nil {
// 		fmt.Printf("Error getting verification code: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("Verification code received: %v\n", response.GetCode())
// }
