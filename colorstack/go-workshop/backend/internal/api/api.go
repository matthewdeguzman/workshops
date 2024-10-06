package api

import (
	"awesome-qrcode-generator/internal/db"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

const MAX_UPLOAD_SIZE int64 = 1024 * 1024 // 1 MB
var qrcodePath = filepath.Join("db", "qrcodes")

var qrdb = db.GetInstance()

type MalformedRequest struct {
	Status int
	Msg    string
}

type DeletePayload struct {
	Id string `json:"id"`
}

type UpdatePayload struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func imagePath(id string) string {
	return filepath.Join(qrcodePath, fmt.Sprintf("%s.png", id))
}

func GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Parsing form")

	// 1. Parse request
	err := r.ParseMultipartForm(1024 * 10) // 10 KB
	if err != nil {
		log.Println("Error ocurred during form processing:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	url, title, description := r.FormValue("url"), r.FormValue("title"), r.FormValue("description")
	log.Println("URL:", url)
	log.Println("Title:", title)
	log.Println("Description:", description)

	// 2. Validate Input
	if url == "" {
		log.Println("Empty URL provided")
		http.Error(w, "Must provide a URL", http.StatusBadRequest)
		return
	}

	// 3. Add QR code to DB
	log.Println("Saving QR code to db")
	qr, err := qrdb.Write(url, title, description)
	if err != nil {
		log.Println("Unable to write to db:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 4. Write QR code to database
	path := imagePath(qr.Id)
	log.Println("Generating QR code to file:", path)
	err = qrcode.WriteFile(url, qrcode.Medium, 512, path)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to generate QR Code", http.StatusInternalServerError)
		return
	}

	// 5. Return id of the code
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(qr)
	if err != nil {
		log.Println("Unable to encode response:", err)
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully generated code")
}

func DeleteQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving id from url")

	// 1. Parse input
	id := r.PathValue("id")

	// 2. Delete id from database
	log.Println("Deleting QR code from database")
	err := qrdb.Delete(id)
	if err != nil {
		if err == db.CodeNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		log.Println("Unsuccessful deletion")
		return
	}

	// 3. Delete id from storage
	log.Println("Deleting QR code from images")
	path := imagePath(id)
	err = os.Remove(path)
	if err != nil {
		log.Printf("Unable to remove code from path '%s'\n", path)
		http.Error(w, "Unable to delete QR code", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully deleted code")
}

func UpdateQRCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Decoding body")

	// 1. Parse input
	err := r.ParseMultipartForm(1024 * 10)
	if err != nil {
		log.Println("Unable to parse form:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	id, title, description := r.FormValue("id"), r.FormValue("title"), r.FormValue("description")

	// 2. Validate input
	if id == "" {
		log.Println("id not provided")
		http.Error(w, "id not provided", http.StatusBadRequest)
		return
	}

	// 3. Update resource in DB
	log.Println("Updating QR code")
	err = qrdb.Update(id, title, description)
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
	log.Println("Fetching Qr Code")

	// 1. Parse input
	id := r.URL.Query().Get("id")
	all := r.URL.Query().Get("all")

	// 2. Validate input
	if all != "" && id != "" {
		http.Error(w, "Request must specify either 'id' or 'all' but not both", http.StatusBadRequest)
		return
	}

	// 3a. Return all QR codes
	if all == "true" {
		log.Println("Fetching all QR codes")
		codes, err := qrdb.GetAll()
		if err != nil {
			log.Println("Unable to fetch QR codes:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Println("Fetched all codes!")

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(codes)
		if err != nil {
			log.Println("Unable to encode response:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// 3b. Return QR code associated with the id
	qr, err := qrdb.GetById(id)
	if err != nil {
		if err == db.CodeNotFound {
			log.Printf("Code with id '%s' not found\n", id)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	log.Println("Code found!")

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(qr)
	if err != nil {
		log.Println("Unable to encode response:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully encoded and code")
}

func GetQRCodeImage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// 1. Search database for QR Code
	_, err := qrdb.GetById(id)
	if err != nil {
		msg := ""
		if err == db.CodeNotFound {
			msg = fmt.Sprintf("Image with id '%s' not found", id)
			http.Error(w, msg, http.StatusNotFound)
		} else {
			msg = "Unable to fetch image"
			http.Error(w, msg, http.StatusInternalServerError)
		}
		log.Println(msg)
		return
	}

	// 2. Serve image
	http.ServeFile(w, r, imagePath(id))
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
