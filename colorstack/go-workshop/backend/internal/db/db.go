package db

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"os"
	"sync"
)

type DB struct {
	qrcodes map[string]*QRCode
}

// TODO: Add title and description and remove size
type QRCode struct {
	ID   string `json:"id"`
	Url  string `json:"url"`
	Size uint32 `json:"size"`
	Path string `json:"path"`
}

var instance *DB
var singletonMtx = &sync.Mutex{}
var dbMtx = &sync.Mutex{}

func (db *DB) init() {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	db.qrcodes = make(map[string]*QRCode)
	file, err := os.ReadFile("db.json")
	if err != nil {
		log.Println("Failed to load db.json")
		log.Println(err)
		panic("Failed to load db.json")
	}

	json.Unmarshal(file, &db.qrcodes)
}

func (db *DB) save() {
	dbMtx.Lock()
	defer dbMtx.Unlock()
	json, err := json.Marshal(db.qrcodes)
	if err != nil {
		log.Println("Failed to save db.json")
		log.Println(err)
		panic("Failed to save db.json")
	}
	os.WriteFile("db.json", json, os.ModePerm)
}

func (db *DB) GetAll() ([]*QRCode, error) {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	var qrcodes []*QRCode
	for _, qr := range db.qrcodes {
		qrcodes = append(qrcodes, &QRCode{ID: qr.ID, Url: qr.Url, Size: qr.Size, Path: qr.Path})
	}

	return qrcodes, nil
}

func (db *DB) GetById(id string) (*QRCode, error) {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	if _, ok := db.qrcodes[id]; !ok {
		return nil, errors.New("QR code not found")
	}

	return &QRCode{ID: db.qrcodes[id].ID, Url: db.qrcodes[id].Url, Size: db.qrcodes[id].Size, Path: db.qrcodes[id].Path}, nil
}

func (db *DB) Delete(id string) error {
	dbMtx.Lock()
	defer dbMtx.Unlock()
	defer db.save()

	if _, ok := db.qrcodes[id]; !ok {
		return errors.New("QR code not found")
	}

	delete(db.qrcodes, id)
	return nil
}

func (db *DB) Update(qr QRCode) error {

	dbMtx.Lock()
	defer dbMtx.Unlock()
	defer db.save()

	if _, ok := db.qrcodes[qr.ID]; !ok {
		return errors.New("QR code not found")
	}

	db.qrcodes[qr.ID] = &qr
	return nil
}

func (db *DB) Write(url string, size uint32, path string) (QRCode, error) {

	id := uuid.New().String()
	qr := &QRCode{
		ID:   id,
		Url:  url,
		Size: size,
		Path: path,
	}

	dbMtx.Lock()
	defer dbMtx.Unlock()
	defer db.save()

	db.qrcodes[id] = qr

	return *qr, nil
}

func GetInstance() *DB {
	if instance == nil {
		// Prevent data races :)
		singletonMtx.Lock()
		defer singletonMtx.Unlock()
		if instance == nil {
			instance = &DB{}
			instance.init()
		}
	}
	return instance
}
