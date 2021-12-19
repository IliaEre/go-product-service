package aws

import (
	"aws-school-service/pkg/domain"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DinamoDbService interface {
	FindAll() (*dynamodb.ScanOutput, error)
	FindOne(id string) (*dynamodb.ScanOutput, error)
	Create(product domain.Product) error
}

func CreateConnection() *dynamodb.DynamoDB {
	log.Println("init dynamodb")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}
