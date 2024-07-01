package main

import (
	"GateServer/rpc"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	pool := rpc.NewRPConPool(5, "localhost", "50051") // Adjust host and port accordingly
	client := rpc.NewVerifyGrpcClient(pool)
	dao, err := NewMysqlDao()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	router.GET("/get_test", func(c *gin.Context) {
		// 获取所有查询参数
		queryParams := c.Request.URL.Query()

		// 创建一个map来存储键值对
		queryMap := make(map[string]string)

		for key, values := range queryParams {
			// 由于每个键可能有多个值，我们这里只取第一个值
			queryMap[key] = values[0]
		}

		c.JSON(http.StatusOK, queryMap)
	})

	router.POST("/get_varifycode", func(c *gin.Context) {
		var requestBody map[string]string
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": Error_Json})
			log.Fatal(err.Error())
			return
		}

		email, exists := requestBody["email"]

		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": Error_Json})
			return
		}
		//执行grpc
		client.GetVarifyCode(email)
		// Do something with the email
		c.JSON(http.StatusOK, gin.H{"email": email, "error": Success})
	})

	router.POST("/user_register", func(c *gin.Context) {
		var requestBody map[string]string
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": Error_Json})
			log.Println("Failed to parse JSON data!", err)
			return
		}

		email, emailExists := requestBody["email"]
		name, nameExists := requestBody["user"]
		pwd, pwdExists := requestBody["passwd"]
		confirm, confirmExists := requestBody["confirm"]
		varifycode, varifycodeExists := requestBody["varifycode"]

		if !emailExists || !nameExists || !pwdExists || !confirmExists || !varifycodeExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": Error_Json})
			return
		}

		if pwd != confirm {
			log.Println("password err")
			c.JSON(http.StatusOK, gin.H{"error": PasswdErr})
			return
		}

		// Check the verification code in Redis

		bGetVarify, err := getRedis(CodePrefix + email)
		if err != nil {
			log.Println("get varify code expired")
			c.JSON(http.StatusOK, gin.H{"error": VarifyExpired})
			return
		}

		if varifycode != bGetVarify {
			log.Println("varify code error", varifycode, "!=", bGetVarify)
			c.JSON(http.StatusOK, gin.H{"error": VarifyCodeErr})
			return
		}

		// Check the database to see if the user exists
		uid, err := dao.RegUser(name, email, pwd)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		if uid == 0 || uid == -1 {
			log.Println("user or email exist")
			c.JSON(http.StatusOK, gin.H{"error": UserExist})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"error":      0,
			"uid":        uid,
			"email":      email,
			"user":       name,
			"passwd":     pwd,
			"confirm":    confirm,
			"varifycode": varifycode,
		})
	})

	router.POST("/reset_pwd", func(c *gin.Context) {
		var requestBody map[string]string
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": Error_Json})
			log.Println("Failed to parse JSON data!", err)
			return
		}

		email, emailExists := requestBody["email"]
		name, nameExists := requestBody["user"]
		pwd, pwdExists := requestBody["passwd"]
		varifycode, varifycodeExists := requestBody["varifycode"]

		if !emailExists || !nameExists || !pwdExists || !varifycodeExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": Error_Json})
			return
		}

		// Check the verification code in Redis
		bGetVarify, err := getRedis(CodePrefix + email)
		if err != nil {
			log.Println("get varify code expired")
			c.JSON(http.StatusOK, gin.H{"error": VarifyExpired})
			return
		}

		if varifycode != bGetVarify {
			log.Println("varify code error", varifycode, "!=", bGetVarify)
			c.JSON(http.StatusOK, gin.H{"error": VarifyCodeErr})
			return
		}

		// Check if the username and email match in the database
		emailValid, err := dao.CheckEmail(name, email)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		if !emailValid {
			log.Println("user email not match")
			c.JSON(http.StatusOK, gin.H{"error": EmailNotMatch})
			return
		}
		pwd = xorString(pwd)
		// Update the password to the new password
		err = dao.UpdatePwd(name, pwd)
		if err != nil {
			log.Println("update pwd failed")
			c.JSON(http.StatusOK, gin.H{"error": PasswdUpFailed})
			return
		}

		log.Println("succeed to update password", pwd)
		c.JSON(http.StatusOK, gin.H{
			"error":      0,
			"email":      email,
			"user":       name,
			"passwd":     pwd,
			"varifycode": varifycode,
		})
	})

	router.Run(":8080")
}

// xor performs XOR operation on the input string with the given key
func xorString(input string) string {
	length := len(input)
	xorCode := length % 255
	result := make([]rune, length)

	for i, char := range input {
		result[i] = char ^ rune(xorCode)
	}

	return string(result)
}
