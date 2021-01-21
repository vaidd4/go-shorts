package app

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	modeAllRUserW    os.FileMode = 0644
	shortsDBFile     string      = "db/shorts.csv"
	addressReadLimit int64       = 1024 * 4 // n chars * utf8 max size
)

//OpenMode file opening mode
type OpenMode int

const (
	rdonly OpenMode = iota
	wronly
	rdwr
)

//OpenShortsDB open file db
func OpenShortsDB(m OpenMode) (*os.File, error) {
	switch m {
	case rdonly:
		return os.Open(shortsDBFile)
	case wronly:
		return os.OpenFile(shortsDBFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, modeAllRUserW)
	case rdwr:
		return os.OpenFile(shortsDBFile, os.O_RDWR|os.O_CREATE, modeAllRUserW)
	default:
		return nil, errors.New("Unknown error")
	}
}

// GET http://localhost:8080/shorts/

func getShorts(w http.ResponseWriter, r *http.Request) {
	// Open DB
	file, err := OpenShortsDB(rdonly)
	if err != nil {
		log.Println("OpenShortsDB() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get document reader
	data := csv.NewReader(file)
	// Read data
	records, err := data.ReadAll()
	// Return data
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(records); err != nil {
		log.Println("json.Encode() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Println("Shorts:", records)
}

// POST http://localhost:8080/shorts/

func createShort(w http.ResponseWriter, r *http.Request) {
	// Open DB file
	file, err := OpenShortsDB(wronly)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read Body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, addressReadLimit))
	if err != nil {
		log.Println("ioutil.ReadAll() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	u, err := url.ParseRequestURI(string(body))
	if err != nil {
		log.Println("url.Parse(body) error:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Generate name
	id, err := nanoid.New(5)
	if err != nil {
		log.Println("nanoid error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get document writer
	data := csv.NewWriter(file)
	// Buffer data
	if err := data.Write([]string{id, u.String()}); err != nil {
		log.Println("data.Write() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Write data
	data.Flush()
	if err := data.Error(); err != nil {
		log.Println("data.Error():", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Write to file synchronously
	if err := file.Sync(); err != nil {
		log.Println("file.Sync() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Short:", id, u.String())
	w.Write([]byte(id))
}

func removeShort(w http.ResponseWriter, r *http.Request, head string) {
	// Open DB file
	file, err := OpenShortsDB(rdwr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get document reader
	reader := csv.NewReader(file)
	// Read data
	records, err := reader.ReadAll()

	var filteredRecords [][]string
	for _, record := range records {
		if record[0] != head {
			filteredRecords = append(filteredRecords, record)
		}
	}

	if len(filteredRecords) == len(records) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Erase file
	if err := file.Truncate(0); err != nil {
		log.Println("file.Truncate() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		log.Println("file.Seek() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get document writer
	writer := csv.NewWriter(file)

	// Buffer data
	if err := writer.WriteAll(filteredRecords); err != nil {
		log.Println("writer.WriteAll() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write data
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Println("data.Error():", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write to file synchronously
	if err := file.Sync(); err != nil {
		log.Println("file.Sync() error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Deleted:", head)
	w.WriteHeader(http.StatusOK)
}
