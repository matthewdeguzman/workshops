package db

type QRCodeConfig struct {
	Url         string `json:"url"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type QRCode struct {
	ID string `json:"id"`
	QRCodeConfig
}
