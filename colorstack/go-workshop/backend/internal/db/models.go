package db

type QrCodeData struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type QrCode struct {
	Id string `json:"id"`
	QrCodeData
}
