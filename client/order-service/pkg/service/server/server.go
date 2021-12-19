package service

import (
	"fmt"
	"io/ioutil"
	"log"

	"encoding/json"
	"net/http"

	api "order-service/pkg/api"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	topic = "aws_test_orders_topic"
)

type OrderService struct {
	api.HttpOrderServer
}

func (hps *OrderService) Create(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(requestBody, "Problem with body!")
	}

	sess := session.Must((session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	svc := sns.New(sess)

	message := string(requestBody)
	result, err := svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &topic,
	})

	if err != nil {
		fmt.Println("problem with sns", err)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*result.MessageId)
}

func Handle(w http.ResponseWriter, message string, code int, err error) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(code)
	log.Println(message, err)
	w.Write([]byte(message))
}
