package main

import (
	"log"
	"net/http"
	"path"
	"strings"
)

func main() {
	app := http.NewServeMux()
	app.HandleFunc("/", appHandler)
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", app))
}

//ShiftPath shift for each segments of path
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
