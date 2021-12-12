package server

import (
	"log"
	"net/http"
)

func HandleException(w http.ResponseWriter, message string, code int, err error) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(code)
	log.Println(message, err)
	w.Write([]byte(message))
}
