package main

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

// session
var svc *dynamodb.DynamoDB
var tableName = "Products"

func findAll(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Products)
}

func findOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("Try to find product with id: ", id)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(id),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		msg := "Could not find '" + id + "'"
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
	} else {
		product := Product{}
		err = dynamodbattribute.UnmarshalMap(result.Item, &product)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)

	}
}

func create(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(requestBody, "Problem with body!")
	}

	var product Product
	json.Unmarshal(requestBody, &product)

	av, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		log.Fatalf("Got error marshalling new product item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

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

	log.Println("init dynamodb")
	// config inside EC2
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = dynamodb.New(sess)
	log.Println("end of init")
}

func main() {
	log.Println("init hander and start server")
	handleRequests()
}
