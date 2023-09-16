package main

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// New initializes a new mail.Dialer instace
func New(smtp SMTP) Mailer {

	host, port, username, password, sender := smtp.Host, smtp.Port, smtp.Username, smtp.Password, smtp.Sender

	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send() takes the recipient email address, file name containing the templates, and any dynamic data for the templates
// func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
func (m Mailer) Send(data MailData) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+data.TemplateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", data.RecepientEmail)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	log.Println("=> Mail sent!")
	return nil
}
