package api

import "net/http"

type ExceptionHandler interface {
	Handle(w http.ResponseWriter, message string, code int, err error)
}
