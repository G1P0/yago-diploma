package handlers

import (
	"net/http"
)

var fs = http.FileServer(http.Dir("web"))

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fs.ServeHTTP(w, r)
}
