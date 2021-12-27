package service

import (
	"aws-school-service/pkg/service/product"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var (
	apiName = "aws-product-service"
)

type Server struct {
	ps *product.ProductService
	ms *http.Server
}

func NewServer(s *product.ProductService) *Server {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(otelmux.Middleware(apiName))
	xraySegment := xray.NewFixedSegmentNamer(apiName)

	myRouter.Handle("/api/health", xray.Handler(xraySegment,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
		})))

	myRouter.Handle("/product/all", xray.Handler(xraySegment, http.HandlerFunc(s.FindAll)))
	myRouter.Handle("/product/{id}", xray.Handler(xraySegment, http.HandlerFunc(s.FindOne)))
	myRouter.Handle("/product/{id}/add", xray.Handler(xraySegment, http.HandlerFunc(s.Create)))

	srv := &http.Server{
		Handler:      myRouter,
		Addr:         "127.0.0.1:8083",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server will be started with address: %s", srv.Addr)
	return &Server{s, srv}
}

func (s *Server) Run() {
	go func() {
		if err := s.ms.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

func (s *Server) Shutdown(context context.Context) {
	s.ms.Shutdown(context)
}
