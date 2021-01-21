package main

import (
	"log"
	"net/http"

	"github.com/vaidd4/go-shorts/app"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.RootHandler)
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
