package api

import (
	"awesome-qrcode-generator/internal/db"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	"io"
	"log"
	"net/http"
	"strings"
)

const MAX_UPLOAD_SIZE int64 = 1024 * 1024 // 1 MB
const QR_CODE_PATH = "./db/qrcodes"

var qrdb = db.GetInstance()

type MalformedRequest struct {
	Status int
	Msg    string
}

type DeletePayload struct {
	Id string `json:"id"`
}

type UpdatePayload struct {
	Id string `json:"id"`
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating QR code")

	var payload db.QRCodeConfig

	err := DecodeJSONBody(w, r, &payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("Unable to decode request body")
		return
	}

	filePath := fmt.Sprintf("%s/%s.png", QR_CODE_PATH, uuid.New().String())
	log.Println("Writing to file:", filePath)
	err = qrcode.WriteFile(payload.Url, qrcode.Medium, 256, filePath)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("Saving QR code to db")
	qr, err := qrdb.Write(payload)
	if err != nil {
		log.Println("Unable to write to db:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.ResponseWriter.WriteHeader(w, http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qr)
	log.Println("Successfully generated code")
}

func DeleteQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Decoding body")
	var payload DeletePayload
	err := DecodeJSONBody(w, r, &payload)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("Deleting QR code")
	err = qrdb.Delete(payload.Id)
	if err != nil {
		if err == db.CodeNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		log.Println("Unsuccessful deletion")
		return
	}
	log.Println("Successfully deleted code")
}

func UpdateQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Decoding body")
	var payload UpdatePayload
	err := DecodeJSONBody(w, r, &payload)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("Updating QR code")
	err = qrdb.Delete(payload.Id)
	if err != nil {
		if err == db.CodeNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		log.Println("Unsuccessful update")
		return
	}
	log.Println("Successfully update code")
}

func GetQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting QR code")
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return err
		}
	}

	// Body may contain additional json. Call decode again to check for this.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}
