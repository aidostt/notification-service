package server

import (
	"notification-service/internal/delivery"
)

func (s *Server) RegisterServers(h *delivery.Handler) {
	// auth.RegisterAuthServer(s.GrpcServer, h)
	// user.RegisterUserServer(s.GrpcServer, h)

}
