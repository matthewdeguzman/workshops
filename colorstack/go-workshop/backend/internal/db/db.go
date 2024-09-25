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

const DB_LOCATION = "./db/db.json"

var instance *DB
var singletonMtx = &sync.Mutex{}
var dbMtx = &sync.Mutex{}

func (db *DB) init() {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	db.qrcodes = make(map[string]*QRCode)
	file, err := os.ReadFile(DB_LOCATION)
	if err != nil {
		log.Println("[DB] Failed to load db.json")
		log.Println(err)
		panic("[DB] Failed to load db.json")
	}

	json.Unmarshal(file, &db.qrcodes)
}

func (db *DB) save() {
	log.Println("[DB] Saving db.json")

	dbMtx.Lock()
	defer dbMtx.Unlock()

	json, err := json.Marshal(db.qrcodes)
	if err != nil {
		log.Println("[DB] Failed to parse qr codes")
		log.Println(err)
	}

	os.WriteFile(DB_LOCATION, json, os.ModePerm)
	if err != nil {
		log.Println("[DB] Failed to save db.json")
		log.Println(err)
	}
	log.Println("[DB] Finished writing to db.json")
}

func (db *DB) saveAndUnlock() {
	dbMtx.Unlock()
	db.save()
}

func (db *DB) GetAll() ([]*QRCode, error) {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	var qrcodes []*QRCode
	for _, qr := range db.qrcodes {
		qrcodes = append(qrcodes, qr)
	}

	return qrcodes, nil
}

func (db *DB) GetById(id string) (*QRCode, error) {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	if _, ok := db.qrcodes[id]; !ok {
		return nil, errors.New("QR code not found")
	}

	return &QRCode{
		ID:           db.qrcodes[id].ID,
		QRCodeConfig: db.qrcodes[id].QRCodeConfig}, nil
}

func (db *DB) Delete(id string) error {
	dbMtx.Lock()
	defer db.saveAndUnlock()

	if _, ok := db.qrcodes[id]; !ok {
		return errors.New("QR code not found")
	}

	delete(db.qrcodes, id)
	return nil
}

func (db *DB) Update(qr QRCode) error {

	dbMtx.Lock()
	defer db.saveAndUnlock()

	if _, ok := db.qrcodes[qr.ID]; !ok {
		return errors.New("QR code not found")
	}

	db.qrcodes[qr.ID] = &qr
	return nil
}

func (db *DB) Write(qrcode QRCodeConfig) (QRCode, error) {

	id := uuid.New().String()
	qr := &QRCode{ID: id, QRCodeConfig: qrcode}

	dbMtx.Lock()
	defer db.saveAndUnlock()

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
