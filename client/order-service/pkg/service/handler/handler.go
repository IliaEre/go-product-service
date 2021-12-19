package handler

import (
	service "order-service/pkg/service/server"

	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func InitHandlers(s *service.OrderService) *http.Server {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(otelmux.Middleware("my-api"))
	xraySegment := xray.NewFixedSegmentNamer("aws-go-service")

	myRouter.Handle("/api/health", xray.Handler(xraySegment,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
		})))
	myRouter.Handle("/order", xray.Handler(xraySegment, http.HandlerFunc(s.Create)))

	srv := &http.Server{
		Handler:      myRouter,
		Addr:         "127.0.0.1:8085",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server will be started with address: %s", srv.Addr)
	return srv
}
