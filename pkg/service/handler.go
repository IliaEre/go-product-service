package service

import (
	// d "aws-school-service/pkg/domain"

	"aws-school-service/pkg/api"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var svc *dynamodb.DynamoDB
var tableName = "Products"

// TODO: remove db
func (h *api.HttpProductServer) FindAll(w http.ResponseWriter, r *http.Request) {
	filt := expression.Name("Id").AttributeExists()

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		// api.(w, "problem...", http.StatusInternalServerError, err)
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	result, err := svc.Scan(input)
	if err != nil {
		hdl.HandleException(w, "problem...", http.StatusInternalServerError, err)
	}

	if len(result.Items) == 0 {
		hdl.HandleException(w, "Not found", http.StatusNotFound, err)
	}

	products := []d.Product{}
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &products); err != nil {
		message := "Got error unmarshalling:"
		hdl.HandleException(w, message, http.StatusInternalServerError, err)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// TODO: remove db
func (h *api.HttpProductServer) FindOne(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
	// log.Println("Try to find product with id: ", id)

	// filt := expression.Name("Id").Equal(expression.Value(id))
	// proj := expression.NamesList(expression.Name("Id"), expression.Name("Name"), expression.Name("Description"))
	// expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	// if err != nil {
	// 	message := "Got error building expression:"
	// 	hdl.HandleException(w, message, http.StatusInternalServerError, err)
	// }

	// params := &dynamodb.ScanInput{
	// 	ExpressionAttributeNames:  expr.Names(),
	// 	ExpressionAttributeValues: expr.Values(),
	// 	FilterExpression:          expr.Filter(),
	// 	ProjectionExpression:      expr.Projection(),
	// 	TableName:                 aws.String(tableName),
	// }

	// result, err := svc.Scan(params)
	// if err != nil {
	// 	message := "Query API call failed:"
	// 	hdl.HandleException(w, message, http.StatusInternalServerError, err)
	// }

	// w.Header().Add("Content-Type", "application/json")
	// var resultSet []d.Product

	// for _, i := range result.Items {
	// 	product := d.Product{}
	// 	err = dynamodbattribute.UnmarshalMap(i, &product)

	// 	if err != nil {
	// 		message := "Got error unmarshalling:"
	// 		hdl.HandleException(w, message, http.StatusInternalServerError, err)
	// 	}
	// 	resultSet = append(resultSet, product)
	// }

	// json.NewEncoder(w).Encode(resultSet)
}

// TODO: remove db
func (h *api.HttpProductServer) Create(w http.ResponseWriter, r *http.Request) {
	// requestBody, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Println(requestBody, "Problem with body!")
	// }

	// var product d.Product
	// json.Unmarshal(requestBody, &product)

	// av, err := dynamodbattribute.MarshalMap(product)
	// if err != nil {
	// 	message := "Got error marshalling new product item:"
	// 	hdl.HandleException(w, message, http.StatusInternalServerError, err)
	// }
	// log.Println("Product:", product)

	// input := &dynamodb.PutItemInput{
	// 	Item:      av,
	// 	TableName: aws.String(tableName),
	// }

	// _, err = svc.PutItem(input)
	// if err != nil {
	// 	message := "Got error calling PutItem:"
	// 	hdl.HandleException(w, message, http.StatusInternalServerError, err)
	// }

	// w.Header().Add("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(product)
}
