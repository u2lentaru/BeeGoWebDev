package controllers

import (
	"io"
	"net/http"
)

// GetPost3 gep post id=3
func GetPost3(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id")
	if key == "3" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "resp": {"post": 3}}`)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status": 500, "err": "db_error"}`)
	}
}
