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

type QrCodeDb struct {
	Path string `json:"path"`
	QrCode
}
