package repository

import "github.com/aws/aws-sdk-go/service/dynamodb"

var svc *dynamodb.DynamoDB
var tableName = "Products"
