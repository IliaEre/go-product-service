package api

import (
	"net/http"
)

type HttpOrderServer interface {
	Create(w http.ResponseWriter, r *http.Request)
}
