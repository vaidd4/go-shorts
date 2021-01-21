package app

import (
	"encoding/csv"
	"net/http"
)

func redirect(w http.ResponseWriter, r *http.Request, head string) {
	// Get db
	file, err := OpenShortsDB(rdonly)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get document reader
	data := csv.NewReader(file)
	// Read data
	records, err := data.ReadAll()

	// Find short
	var redirect string
	for _, v := range records {
		if v[0] == head {
			redirect = v[1]
			break
		}
	}
	if redirect == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	// Redirect
	http.Redirect(w, r, redirect, http.StatusMovedPermanently)
}
