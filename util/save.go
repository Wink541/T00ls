package util

import (
	"fmt"
	"net/smtp"
)

type SMTPServer struct {
	Host string
	Port string
}

type Sender struct {
	SmtpServer SMTPServer
	Email      string
	Password   string
}

type SendInfo struct {
	Sender
	Recipient []string
	Message   []byte
}

func SendMail(info SendInfo) {
	err := smtp.SendMail(
		info.SmtpServer.Host+":"+info.SmtpServer.Port,
		smtp.PlainAuth("", info.Email, info.Password, info.SmtpServer.Host),
		info.Email,
		info.Recipient,
		info.Message)
	if err != nil {
		fmt.Println(err)
	}
}
