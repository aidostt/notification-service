package dialog

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"notification-service/pkg/logger"
)

type DialogService interface {
	NewConnection(string) (*grpc.ClientConn, error)
}

type Dialog struct {
	Addresses Addresses
	authority string
}
type Addresses struct {
	QRs          string
	Users        string
	Reservations string
}

func NewDialog(authority, qrs, reservations, users string) *Dialog {
	return &Dialog{authority: authority, Addresses: Addresses{
		QRs:          qrs,
		Reservations: reservations,
		Users:        users,
	}}
}

func (d *Dialog) NewConnection(address string) (*grpc.ClientConn, error) {
	//cert, err := tls.LoadX509KeyPair("path/to/server.crt", "path/to/server.key")
	//if err != nil {
	//	return nil, err
	//}
	//tlsConfig := &tls.Config{
	//	Certificates: []tls.Certificate{cert},
	//	ServerName:   d.authority,
	//}
	//creds := credentials.NewTLS(tlsConfig)
	//conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Failed to connect: %v", err)
		conn.Close()
		return nil, err
	}
	return conn, nil
}
