package service

import (
	"bytes"
	"embed"
	"html/template"
	"notification-service/pkg/dialog"
	"time"

	gomail "github.com/go-mail/mail"
)

//go:embed templates
var templateFS embed.FS

const (
	templatesPath = "templates/"
)

type MailerService struct {
	dialer *gomail.Dialer
	sender string
	dialog *dialog.Dialog
}

func NewMailerService(host string, port int, username, password, sender string, grpcDialog *dialog.Dialog) *MailerService {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second
	return &MailerService{
		dialer: dialer,
		sender: sender,
		dialog: grpcDialog,
	}
}

func (s *MailerService) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, templatesPath+templateFile)
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

	msg := gomail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", s.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		var err error
		err = s.dialer.DialAndSend(msg)
		errChan <- err
		return
	}()

	// Wait for the go routine to send an error or nil.
	return <-errChan
}
