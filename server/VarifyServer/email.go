package main

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

// 配置模块中的邮箱用户名和授权码

// 创建发送邮件的代理
func createTransport() *email.Email {
	e := email.NewEmail()
	e.From = EmailUser
	return e
}

// 发送邮件的函数
func sendMail(mailOptions *email.Email) (string, error) {
	e := createTransport()
	e.To = mailOptions.To
	e.Subject = mailOptions.Subject
	e.Text = mailOptions.Text
	e.HTML = mailOptions.HTML

	auth := smtp.PlainAuth("", EmailUser, EmailPass, "smtp.qq.com")
	err := e.Send("smtp.qq.com:587", auth)
	if err != nil {
		return "", err
	}

	return "邮件已成功发送", nil
}

// func main() {
// 	// 示例：发送邮件
// 	mailOptions := &email.Email{
// 		To:      []string{"1820737440@qq.com"},
// 		Subject: "测试邮件",
// 		Text:    []byte("这是测试邮件内容"),
// 	}

// 	response, err := sendMail(mailOptions)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Println(response)
// 	}
// }
