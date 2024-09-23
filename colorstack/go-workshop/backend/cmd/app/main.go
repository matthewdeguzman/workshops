package main

import (
	"awesome-qrcode-generator/internal/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /api/generate", api.GenerateQRCode)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
