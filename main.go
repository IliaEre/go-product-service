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

	httpSwagger "github.com/swaggo/http-swagger"
)

// GetProducts godoc
// @Summary Get details of products
// @Description Get details of products
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {array} Order
// @Router /product/all [get]
func findAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Products)
}

// GetProducts godoc
// @Summary Get details of product
// @Description Get details of product
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {array} Order
// @Router /product/{id} [get]
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

// CreateOrder godoc
// @Summary Create a new product
// @Description Create a new product with the input paylod
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body Order true "Create product"
// @Success 200 {object} Order
// @Router /product [post]
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

// CreateOrder godoc
// @Summary Update a new product
// @Description Update a selected product with the input paylod
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body Order true "Create order"
// @Success 200 {object} Order
// @Router /product/{id}/add [put]
func update(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(requestBody, "Problem with body!")
	}

	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		log.Println("problem with id...", err)
	}

	w.Header().Add("Content-Type", "application/json")
	var product Product
	json.Unmarshal(requestBody, &product)

	for index, doc := range Products {
		if int64(doc.Id) == key {
			Products[index] = product
			json.NewEncoder(w).Encode(product)
		}
	}
}

// @title Orders API
// @version 1.0
// @description This is a sample serice for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8082
// @BasePath /
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(otelmux.Middleware("my-api"))
	xraySegment := xray.NewFixedSegmentNamer("aws-go-service")

	myRouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	myRouter.Handle("/api/health", xray.Handler(xraySegment,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
		})))

	myRouter.Handle("/product", xray.Handler(xraySegment, http.HandlerFunc(create)))
	myRouter.Handle("/product/all", xray.Handler(xraySegment, http.HandlerFunc(findAll)))
	myRouter.Handle("/product/{id}", xray.Handler(xraySegment, http.HandlerFunc(findOne)))
	myRouter.Handle("/product/{id}/add", xray.Handler(xraySegment, http.HandlerFunc(update)))
	http.ListenAndServe(":8082", myRouter)
}

func init() {
	log.Println("init products")
	Products = []Product{
		{Id: 1, Name: "Useful smartphone", Description: "There is a description here..."},
		{Id: 2, Name: "Useful laptop", Description: "khm..."},
	}

	log.Println("init x-ray configuration")
	s, _ := sampling.NewCentralizedStrategyWithFilePath("rules.json")
	xray.Configure(xray.Config{SamplingStrategy: s})
	log.Println("end of init")
}

func main() {
	log.Println("init hander and start server")
	handleRequests()
}
