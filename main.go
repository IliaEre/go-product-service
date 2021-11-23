package main

import (
	"encoding/json"
	"os"
	"strconv"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
)

func findAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	xrayTraceId := getXrayTraceID(trace.SpanFromContext(ctx))

	response := ResponseEntityAll{Products, xrayTraceId}
	json.NewEncoder(w).Encode(response)
}

func findOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		log.Println("problem with id...", err)
	}

	for _, doc := range Products {
		if int64(doc.Id) == key {
			json.NewEncoder(w).Encode(doc)
		}
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(requestBody, "Problem with body!")
	}

	var product Product
	json.Unmarshal(requestBody, &product)
	Products = append(Products, product)
	json.NewEncoder(w).Encode(product)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(otelmux.Middleware("my-api"))
	xraySegment := xray.NewFixedSegmentNamer("aws-go-service")

	myRouter.Handle("/api/health", xray.Handler(xraySegment,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			xrayTraceID := getXrayTraceID(trace.SpanFromContext(ctx))
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "traceId": xrayTraceID})
		})))

	myRouter.Handle("/product/all", xray.Handler(xraySegment, http.HandlerFunc(findAll)))
	myRouter.Handle("/product/{id}", xray.Handler(xraySegment, http.HandlerFunc(findOne)))
	myRouter.Handle("/product/{id}/add", xray.Handler(xraySegment, http.HandlerFunc(create)))
	http.ListenAndServe(":8082", myRouter)
}

func init() {
	log.Println("init products")
	Products = []Product{
		{Id: 1, Name: "Useful smartphone", Description: "There is a description here..."},
		{Id: 2, Name: "Useful laptop", Description: "khm..."},
	}

	log.Println("init x-ray configudaration")
	// ctx := context.Background()
	if os.Getenv("ENVIRONMENT") == "production" {
		ec2.Init()
	}

	s, _ := sampling.NewCentralizedStrategyWithFilePath("rules.json")
	xray.Configure(xray.Config{SamplingStrategy: s})
	log.Println("end of init")
}

func main() {
	log.Println("init hander and start server")
	handleRequests()
}

func getXrayTraceID(span trace.Span) string {
	xrayTraceID := span.SpanContext().TraceID().String()
	result := fmt.Sprintf("1-%s-%s", xrayTraceID[0:8], xrayTraceID[8:])
	return result
}
