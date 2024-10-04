package main

import (
	"awesome-qrcode-generator/internal/api"
	"awesome-qrcode-generator/internal/middleware"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("OPTIONS /", middleware.Chain(
		middleware.Preflight,
		middleware.Time(),
		middleware.Logger()))

	http.HandleFunc("POST /api/generate", middleware.Chain(
		api.GenerateQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth(),
		middleware.Cors()))

	http.HandleFunc("DELETE /api/delete/{id}", middleware.Chain(
		api.DeleteQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth(),
		middleware.Cors()))

	http.HandleFunc("PUT /api/update/", middleware.Chain(
		api.UpdateQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth(),
		middleware.Cors()))

	http.HandleFunc("GET /api/get", middleware.Chain(
		api.GetQRCode,
		middleware.Time(),
		middleware.Logger(),
		middleware.Auth(),
		middleware.Cors()))

	http.HandleFunc("GET /qr/{id}", middleware.Chain(
		api.GetQRCodeImage,
		middleware.Time(),
		middleware.Logger(),
		middleware.Cors()))

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
