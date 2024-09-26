package db

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/teris-io/shortid"
)

type DB struct {
	qrcodes map[string]*QrCodeDb
}

const DB_LOCATION = "./db/db.json"

var sid = shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().Unix()))

var CodeNotFound error = errors.New("QR code not found")
var InternalDbError error = errors.New("Internal Database Error")

var instance *DB
var singletonMtx = &sync.Mutex{}
var dbMtx = &sync.Mutex{}

func (db *DB) init() {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	db.qrcodes = make(map[string]*QrCodeDb)
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
		return
	}

	err = os.WriteFile(DB_LOCATION, json, os.ModePerm)
	if err != nil {
		log.Println("[DB] Failed to save db.json")
		log.Println(err)
		return
	}
	log.Println("[DB] Finished writing to db.json")
}

func (db *DB) saveAndUnlock() {
	dbMtx.Unlock()
	db.save()
}

func (db *DB) GetAll() ([]QrCode, error) {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	var qrcodes []QrCode
	for _, qr := range db.qrcodes {
		qrcodes = append(qrcodes, qr.QrCode)
	}

	return qrcodes, nil
}

func (db *DB) GetById(id string) (QrCodeDb, error) {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	if _, ok := db.qrcodes[id]; !ok {
		return QrCodeDb{}, CodeNotFound
	}

	return *db.qrcodes[id], nil
}

func (db *DB) Delete(id string) error {
	dbMtx.Lock()

	if _, ok := db.qrcodes[id]; !ok {
		dbMtx.Unlock()
		return CodeNotFound
	}
	defer db.saveAndUnlock()

	delete(db.qrcodes, id)
	return nil
}

func (db *DB) Update(id string, title string, description string) error {

	dbMtx.Lock()

	if _, ok := db.qrcodes[id]; !ok {
		dbMtx.Unlock()
		return CodeNotFound
	}
	defer db.saveAndUnlock()

	db.qrcodes[id].QrCode.Title = title
	db.qrcodes[id].QrCode.Description = description
	return nil
}

func (db *DB) Write(url string, title string, description string, path string) (QrCodeDb, error) {
	id, err := sid.Generate()
	if err != nil {
		log.Println("[DB] Unable to create short id:", err)
		return QrCodeDb{}, InternalDbError
	}

	qr := &QrCodeDb{Path: path, QrCode: QrCode{Id: id, QrCodeData: QrCodeData{Url: url, Title: title, Description: description}}}

	dbMtx.Lock()
	defer db.saveAndUnlock()

	db.qrcodes[id] = qr

	return *qr, nil
}

func GetInstance() *DB {
	if instance == nil {
		singletonMtx.Lock()
		defer singletonMtx.Unlock()
		if instance == nil {
			instance = &DB{}
			instance.init()
		}
	}
	return instance
}
