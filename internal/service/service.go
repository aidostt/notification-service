package service

import "notification-service/pkg/dialog"

type Mailer interface {
	Send(recipient, templateFile string, data interface{}) error
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
