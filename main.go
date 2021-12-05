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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var svc *dynamodb.DynamoDB
var tableName = "Products"

func findAll(w http.ResponseWriter, r *http.Request) {
	out, err := svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		message := "Failed to find products:"
		handeException(w, message, http.StatusInternalServerError, err)
	}

	if out != nil {
		w.Header().Add("Content-Type", "application/json")
		var resultSet []Product
		for _, element := range out.Items {
			product := Product{}
			err = dynamodbattribute.UnmarshalMap(element, &product)

			if err != nil {
				message := "Got error unmarshalling:"
				handeException(w, message, http.StatusInternalServerError, err)
			}
			resultSet = append(resultSet, product)
		}

		json.NewEncoder(w).Encode(resultSet)

	} else {
		message := "Could not find products"
		handeException(w, message, http.StatusNotFound, err)
	}
}

func findOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("Try to find product with id: ", id)

	filt := expression.Name("Id").Equal(expression.Value(id))
	proj := expression.NamesList(expression.Name("Id"), expression.Name("Name"), expression.Name("Description"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		message := "Got error building expression:"
		handeException(w, message, http.StatusInternalServerError, err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		message := "Query API call failed:"
		handeException(w, message, http.StatusInternalServerError, err)
	}

	w.Header().Add("Content-Type", "application/json")
	var resultSet []Product

	for _, i := range result.Items {
		product := Product{}
		err = dynamodbattribute.UnmarshalMap(i, &product)

		if err != nil {
			message := "Got error unmarshalling:"
			handeException(w, message, http.StatusInternalServerError, err)
		}
		resultSet = append(resultSet, product)
	}

	json.NewEncoder(w).Encode(resultSet)
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
		message := "Got error marshalling new product item:"
		handeException(w, message, http.StatusInternalServerError, err)
	}
	log.Println("Product:", product)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		message := "Got error calling PutItem:"
		handeException(w, message, http.StatusInternalServerError, err)
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
	log.Println("init x-ray configudaration")
	s, _ := sampling.NewCentralizedStrategyWithFilePath("rules.json")
	xray.Configure(xray.Config{SamplingStrategy: s})

	log.Println("init dynamodb")
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

func handeException(w http.ResponseWriter, message string, code int, err error) {
	w.WriteHeader(code)
	log.Println(message, err)
	w.Write([]byte(message))
}
