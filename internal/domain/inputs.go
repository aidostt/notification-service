package domain

type ContentInput struct {
	Content string
}

type QRCodeMailInput struct {
	QRCodeBase64 string
	User         UserInfo
	Restaurant   RestaurantInfo
}

type UserInfo struct {
	Name    string
	Surname string
	Phone   string
	Email   string
}

type RestaurantInfo struct {
	Name            string
	Address         string
	Contact         string
	Table           int32
	ReservationTime string
}
