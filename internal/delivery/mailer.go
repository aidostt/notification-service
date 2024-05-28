package delivery

import (
	"context"
	proto "github.com/aidostt/protos/gen/go/reservista/mailer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"notification-service/internal/domain"
	"notification-service/pkg/logger"
)

func (h *Handler) SendWelcome(ctx context.Context, input *proto.ContentInput) (*proto.StatusResponse, error) {
	if input.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if input.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	contentInput := domain.ContentInput{
		Content: input.GetContent(),
	}

	err := h.Service.Mailer.Send(input.Email, "welcome_page.html", nil, contentInput)
	if err != nil {
		logger.Error(err)
		return &proto.StatusResponse{Status: false}, status.Error(codes.Internal, "failed to send message: "+err.Error())
	}
	return &proto.StatusResponse{Status: true}, nil
}

func (h *Handler) SendQR(ctx context.Context, input *proto.QRInput) (*proto.StatusResponse, error) {
	if input.GetReservationID() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	if input.GetUserID() == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation id is required")
	}
	if input.GetQRUrlBase() == "" {
		return nil, status.Error(codes.InvalidArgument, "url base is required")
	}

	err := h.Service.Mailer.SendQR("qr_code_email_template.html", ctx, input.GetUserID(), input.GetReservationID(), input.GetQRUrlBase())
	if err != nil {
		logger.Error(err)
		return nil, status.Error(codes.Internal, "failed to send message: "+err.Error())
	}
	return &proto.StatusResponse{Status: true}, nil
}

func (h *Handler) SendAuthCode(ctx context.Context, input *proto.ContentInput) (*proto.StatusResponse, error) {
	if input.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if input.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	contentInput := domain.ContentInput{
		Content: input.GetContent(),
	}

	err := h.Service.Mailer.Send(input.Email, "verification_code.html", nil, contentInput)
	if err != nil {
		logger.Error(err)
		return &proto.StatusResponse{Status: false}, status.Error(codes.Internal, "failed to send message: "+err.Error())
	}
	return &proto.StatusResponse{Status: true}, nil
}

func (h *Handler) SendReminder(ctx context.Context, input *proto.EmailInput) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{Status: true}, nil
	// getUser
	// getReservation.Time
}
func (h *Handler) SendAdvert(ctx context.Context, input *proto.EmailInput) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{Status: true}, nil
	// getUser
	// advert?
}
func (h *Handler) SendResetCode(ctx context.Context, input *proto.EmailInput) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{Status: true}, nil
	// getUser
	// advert?
}
