package app

import (
	"log"
	"net/http"

	"github.com/vaidd4/go-shorts/util"
)

//RootHandler Main handler for the app
func RootHandler(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = util.ShiftPath(r.URL.Path)
	switch head {
	case "":
		interfaceHandler(w, r)
	case "shorts":
		shortsHandler(w, r)
	default:
		redirectHandler(w, r, head)
	}
}

func interfaceHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: C. Create web interface
	log.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		// Return web interface
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func shortsHandler(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = util.ShiftPath(r.URL.Path)
	if head == "" {
		switch r.Method {
		case http.MethodGet:
			getShorts(w, r)
		case http.MethodPost:
			createShort(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			// Return a short.
			http.Error(w, "Not Implemented", http.StatusNotImplemented)
		case http.MethodPut:
			// Update a short.
			http.Error(w, "Not Implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			removeShort(w, r, head)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request, head string) {
	// Handle any url
	switch r.Method {
	case http.MethodGet:
		// Redirect to URL
		redirect(w, r, head)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
