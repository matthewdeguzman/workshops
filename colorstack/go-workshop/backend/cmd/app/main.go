package main

import (
	"awesome-qrcode-generator/internal/api"
	"awesome-qrcode-generator/internal/middleware"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /api/generate", middleware.Chain(
		api.GenerateQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth()))

	http.HandleFunc("DELETE /api/delete", middleware.Chain(
		api.DeleteQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth()))

	http.HandleFunc("PUT /api/update", middleware.Chain(
		api.UpdateQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth()))

	http.HandleFunc("GET /api/get", middleware.Chain(
		api.GetQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth()))

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
