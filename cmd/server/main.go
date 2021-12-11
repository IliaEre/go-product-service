package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
)

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
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	log.Println("init hander and start server")
	srv := initHandlers()
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
