package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	a "aws-school-service/pkg/service/aws"
	srv "aws-school-service/pkg/service/handler"
	service "aws-school-service/pkg/service/server"
)

func init() {
	a.InitXrayConfig()
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	log.Println("init hander and start server")
	se := service.ProductService{}
	srv := srv.InitHandlers(&se)
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
