package api

import (
	"aws-school-service/pkg/domain"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DinamoDbRepository interface {
	FindAll() (*dynamodb.ScanOutput, error)
	FindOne(id string) (*dynamodb.ScanOutput, error)
	Create(product domain.Product) error
}
