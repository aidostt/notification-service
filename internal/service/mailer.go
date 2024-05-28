package service

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	proto_qr "github.com/aidostt/protos/gen/go/reservista/qr"
	proto_reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	proto_user "github.com/aidostt/protos/gen/go/reservista/user"
	"github.com/skip2/go-qrcode"
	"html/template"
	"image/png"
	"notification-service/internal/domain"
	"notification-service/pkg/dialog"
	"notification-service/pkg/logger"
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
	Dialog *dialog.Dialog
}

func NewMailerService(host string, port int, username, password, sender string, grpcDialog *dialog.Dialog) *MailerService {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second
	return &MailerService{
		dialer: dialer,
		sender: sender,
		Dialog: grpcDialog,
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

func (s *MailerService) SendQR(templateFile string, ctx context.Context, userID, reservationID, qrURLBase string) error {
	conn, err := s.Dialog.NewConnection(s.Dialog.Addresses.QRs)
	defer conn.Close()
	if err != nil {
		return err
	}
	qrClient := proto_qr.NewQRClient(conn)
	resp, err := qrClient.Generate(ctx, &proto_qr.GenerateRequest{
		Content: qrURLBase + reservationID,
	})
	if err != nil {
		return err
	}
	qrCodeBase64, err := s.GenerateQRCodeBase64(string(resp.QR))
	if err != nil {
		logger.Error(err)
		return err
	}
	userClient := proto_user.NewUserClient(conn)
	userResponse, err := userClient.GetByID(ctx, &proto_user.GetRequest{
		UserId: userID,
		Email:  "plug",
	})
	if err != nil {
		return err
	}
	reservationClient := proto_reservation.NewReservationClient(conn)
	reservation, err := reservationClient.GetReservation(ctx, &proto_reservation.IDRequest{Id: reservationID})
	if err != nil {
		return err
	}
	user := domain.UserInfo{
		Name:    userResponse.GetName(),
		Surname: userResponse.GetSurname(),
		Phone:   userResponse.GetPhone(),
		Email:   userResponse.GetEmail(),
	}
	restaurant := domain.RestaurantInfo{
		Name:            reservation.Table.Restaurant.GetName(),
		Address:         reservation.Table.Restaurant.GetAddress(),
		Contact:         reservation.Table.Restaurant.GetContact(),
		Table:           reservation.Table.GetTableNumber(),
		ReservationTime: reservation.GetReservationTime(),
	}
	qrInput := domain.QRCodeMailInput{
		QRCodeBase64: qrCodeBase64,
		User:         user,
		Restaurant:   restaurant,
	}
	return s.Send(user.Email, templateFile, qrInput)
}

func (s *MailerService) GenerateQRCodeBase64(data string) (string, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, qr.Image(256))
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Str, nil
}
