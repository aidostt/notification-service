package service

import (
	"bytes"
	"context"
	"embed"
	proto_qr "github.com/aidostt/protos/gen/go/reservista/qr"
	proto_reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	proto_user "github.com/aidostt/protos/gen/go/reservista/user"
	"html/template"
	"notification-service/internal/domain"
	"notification-service/pkg/dialog"
	"os"
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

func (s *MailerService) Send(recipient, templateFile string, qrCode *os.File, data any) error {
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
	if qrCode != nil {
		defer qrCode.Close()
		defer os.Remove(qrCode.Name())
		msg.Embed(qrCode.Name(), gomail.SetHeader(map[string][]string{
			"Content-ID": {"<qrcode>"},
		}))
		go func() {
			defer close(errChan)
			var err error
			err = s.dialer.DialAndSend(msg)
			errChan <- err
			return
		}()
	} else {
		go func() {
			defer close(errChan)
			var err error
			err = s.dialer.DialAndSend(msg)
			errChan <- err
			return
		}()
	}

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
	qrResponse, err := qrClient.Generate(ctx, &proto_qr.GenerateRequest{
		Content: qrURLBase + reservationID,
	})
	if err != nil {
		return err
	}
	conn, err = s.Dialog.NewConnection(s.Dialog.Addresses.Users)
	if err != nil {
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
	conn, err = s.Dialog.NewConnection(s.Dialog.Addresses.Reservations)
	if err != nil {
		return err
	}
	reservationClient := proto_reservation.NewReservationClient(conn)
	reservationResponse, err := reservationClient.GetReservation(ctx, &proto_reservation.IDRequest{Id: reservationID})
	if err != nil {
		return err
	}

	// Extract and parse the reservation time
	reservationTimeStr := reservationResponse.GetReservationTime() // example: "1:00 PM"
	reservationTime, err := time.Parse("3:04 PM", reservationTimeStr)
	if err != nil {
		return err
	}
	formattedTime := reservationTime.Format("15:04 PM")

	// Extract and parse the reservation date
	reservationDate := reservationResponse.GetReservationDate().AsTime() // Converts to time.Time
	formattedDate := reservationDate.Format("Jan 02, 2006")
	user := domain.UserInfo{
		Name:    userResponse.GetName(),
		Surname: userResponse.GetSurname(),
		Phone:   userResponse.GetPhone(),
		Email:   userResponse.GetEmail(),
	}
	restaurant := domain.RestaurantInfo{
		Name:    reservationResponse.Table.Restaurant.GetName(),
		Address: reservationResponse.Table.Restaurant.GetAddress(),
		Phone:   reservationResponse.Table.Restaurant.GetContact(),
		Table:   reservationResponse.Table.GetTableNumber(),
	}
	reservation := domain.ReservationInfo{
		ReservationTime: formattedTime,
		ReservationDate: formattedDate,
		ReservationID:   reservationID,
	}
	data := domain.QRCodeMailInput{
		QR:          qrResponse.GetQR(),
		User:        user,
		Restaurant:  restaurant,
		Reservation: reservation,
	}

	tmpFile, err := os.CreateTemp("", "qrcode-*.png")
	if err != nil {
		return err
	}

	if _, err = tmpFile.Write(data.QR); err != nil {
		return err
	}
	return s.Send(user.Email, templateFile, tmpFile, data)
}
