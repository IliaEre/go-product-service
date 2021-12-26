package repository

import (
	"aws-school-service/pkg/api"
	"aws-school-service/pkg/domain"
	a "aws-school-service/pkg/service/aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var Dynamo *dynamodb.DynamoDB
var TableName = "Products"

type DymanoRepository struct {
	api.DinamoDbRepository
}

func NewRepository() *DymanoRepository {
	Dynamo = a.CreateConnection()
	return &DymanoRepository{}
}

func (r *DymanoRepository) FindAll() (*dynamodb.ScanOutput, error) {
	filt := expression.Name("Id").AttributeExists()
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	result, err := Dynamo.Scan(input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, err
	}

	return result, nil
}

func (r *DymanoRepository) FindOne(id string) (*dynamodb.ScanOutput, error) {
	filt := expression.Name("Id").Equal(expression.Value(id))
	proj := expression.NamesList(expression.Name("Id"), expression.Name("Name"), expression.Name("Description"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(TableName),
	}

	result, err := Dynamo.Scan(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *DymanoRepository) Create(product domain.Product) error {
	av, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableName),
	}

	_, err = Dynamo.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
