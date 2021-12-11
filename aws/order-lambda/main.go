package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"encoding/json"
	"os"
)

var topic = "orderTopic"
var svc *sns.SNS

type OrderEvent struct {
	Name      string `json:"name"`
	ProductId string `json:"productId"`
}

func HandleRequest(ctx context.Context, order OrderEvent) (string, error) {
	json, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshalling order")
		return "", err
	}

	message := string(json)

	result, err := svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &topic,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(*result.MessageId)
	return fmt.Sprintf("Order was created with name: %s!", order.Name), nil
}

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = sns.New(sess)
}

func main() {
	lambda.Start(HandleRequest)
}
