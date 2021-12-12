package repository

import (
	"aws-school-service/pgk/service/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Dynamo *dynamodb.DynamoDB
var TableName = "Products"

func CreateConnection() {
	Dynamo = aws.CreateConnection()
}
