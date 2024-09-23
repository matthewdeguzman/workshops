package api

import (
	"log"
	"net/http"
)

func GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating QR Code")
}
