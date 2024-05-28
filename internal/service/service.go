package service

import (
	"context"
	"notification-service/pkg/dialog"
	"os"
)

type Mailer interface {
	Send(recipient, templateFile string, qrCode *os.File, data interface{}) error
	SendQR(string, context.Context, string, string, string) error
}

type Services struct {
	Mailer Mailer
}

type Dependencies struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
	Dialog   *dialog.Dialog
}

func NewService(deps Dependencies) *Services {
	MailService := NewMailerService(deps.Host, deps.Port, deps.Username, deps.Password, deps.Sender, deps.Dialog)
	return &Services{
		Mailer: MailService,
	}
}
