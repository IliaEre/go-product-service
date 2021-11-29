package main

import (
	"encoding/json"
	"strconv"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func findAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Products)
}

func findOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
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
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(otelmux.Middleware("my-api"))
	xraySegment := xray.NewFixedSegmentNamer("aws-go-service")

	myRouter.Handle("/api/health", xray.Handler(xraySegment,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
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
	s, _ := sampling.NewCentralizedStrategyWithFilePath("rules.json")
	xray.Configure(xray.Config{SamplingStrategy: s})
	log.Println("end of init")
}

func main() {
	log.Println("init hander and start server")
	handleRequests()
}
