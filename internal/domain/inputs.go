package domain

type ContentInput struct {
	Content string
}

type QRCodeMailInput struct {
	QR          []byte
	User        UserInfo
	Restaurant  RestaurantInfo
	Reservation ReservationInfo
}

type UserInfo struct {
	Name    string
	Surname string
	Phone   string
	Email   string
}

type ReservationInfo struct {
	ReservationTime string
	ReservationDate string
	ReservationID   string
}

type RestaurantInfo struct {
	Name    string
	Address string
	Phone   string
	Table   int32
}
