package utils

import (
	"strconv"

	"github.com/vpa/quanlynhahang-backend/config"
	"gopkg.in/gomail.v2"
)

func SendMail(to, subject, body string) error {
	host := config.GetEnv("MAIL_HOST")
	portStr := config.GetEnv("MAIL_PORT")
	username := config.GetEnv("MAIL_USERNAME")
	password := config.GetEnv("MAIL_PASSWORD")
	from := config.GetEnv("MAIL_FROM")

	port, _ := strconv.Atoi(portStr)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		host,
		port,
		username,
		password,
	)

	return d.DialAndSend(m)
}
