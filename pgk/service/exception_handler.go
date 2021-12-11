package service

import (
	"log"
	"net/http"
)

func HandeException(w http.ResponseWriter, message string, code int, err error) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(code)
	log.Println(message, err)
	w.Write([]byte(message))
}
