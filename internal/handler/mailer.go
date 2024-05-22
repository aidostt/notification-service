package handler

import (
	"context"
	proto "github.com/aidostt/protos/gen/go/reservista/mailer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"notification-service/pkg/logger"
)

// in user repo generate verReset code
// getVerCode, getResetCode

func (h *Handler) SendWelcome(ctx context.Context, input *proto.EmailInput) (*proto.StatusResponse, error) {
	if input.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	err := h.Service.Mailer.Send(input.GetEmail(), "user_welcome.html", nil)
	if err != nil {
		logger.Error(err)
		return nil, status.Error(codes.Internal, "failed to send message: "+err.Error())
	}
	return &proto.StatusResponse{Status: true}, nil
}

func (h *Handler) SendQR(ctx context.Context, input *proto.EmailInput) (*proto.StatusResponse, error) {
	if input.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	err := h.Service.Mailer.Send(input.GetEmail(), "qr_code.html", nil)
	if err != nil {
		logger.Error(err)
		return nil, status.Error(codes.Internal, "failed to send message: "+err.Error())
	}
	return &proto.StatusResponse{Status: true}, nil
}

func (h *Handler) SendAuthCode(ctx context.Context, input *proto.EmailInput) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{Status: true}, nil
	// getUser
	// advert?
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
