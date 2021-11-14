package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
	myRouter.HandleFunc("/product/{id}", findOne).Methods("GET")
	myRouter.HandleFunc("/product/{id}/add", create).Methods("PUT")
	myRouter.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	http.ListenAndServe(":8082", myRouter)
}

func init() {
	log.Println("init products")
	Products = []Product{
		{Id: 1, Name: "Useful smartphone", Description: "There is a description here..."},
		{Id: 2, Name: "Useful laptop", Description: "khm..."},
	}
}

func main() {
	log.Println("init hander and start server")
	handleRequests()
}
