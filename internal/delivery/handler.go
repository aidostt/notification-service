package delivery

import (
	proto "github.com/aidostt/protos/gen/go/reservista/mailer"
	"notification-service/internal/service"
)

type Handler struct {
	Service *service.Services
	proto.UnimplementedMailerServer
}

func NewHandler(Service *service.Services) *Handler {
	return &Handler{
		Service: Service,
	}
}
