package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"aws-school-service/pkg/repository"
	a "aws-school-service/pkg/service/aws"
	"aws-school-service/pkg/service/product"
	srv "aws-school-service/pkg/service/server"
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

	repo := repository.NewRepository()
	srvs := product.NewProductService(repo)
	s := srv.NewServer(srvs)

	s.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	s.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
