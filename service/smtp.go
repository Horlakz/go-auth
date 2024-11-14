package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"

	"github.com/horlakz/go-auth/constants"
	"github.com/horlakz/go-auth/dto"
)

var SenderName = "DukiaGold"

type SmtpInterface interface {
	SendWithTemplate(sendEmailDto dto.SendEmailDto) error
}

type smtpService struct {
	Host, Username, Password, From string
	Port                           int
}

func NewSmtpService(env constants.Env) SmtpInterface {
	port, _ := strconv.Atoi(env.SMTP_PORT)

	return &smtpService{
		Host:     env.SMTP_HOST,
		Port:     port,
		Username: env.SMTP_USERNAME,
		Password: env.SMTP_PASSWORD,
		From:     env.FROM_EMAIL,
	}
}

func (e *smtpService) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	msg := []byte("From: " + SenderName + " <" + e.From + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		body)
	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.Host,
	})
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, e.Host)
	if err != nil {
		return err
	}

	// Authenticating
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Setting the sender and recipient
	if err = client.Mail(e.From); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	// Sending the email body
	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return err
	}

	if err = wc.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func (e *smtpService) ParseTemplate(templateFile string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFile, "templates/layout.html")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err = t.ExecuteTemplate(buf, "layout", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (e *smtpService) SendWithTemplate(sendEmailDto dto.SendEmailDto) error {
	body, err := e.ParseTemplate(fmt.Sprintf("templates/%s.html", sendEmailDto.Template), sendEmailDto.Variables)
	if err != nil {
		return err
	}

	return e.Send(sendEmailDto.To, sendEmailDto.Subject, body)
}
