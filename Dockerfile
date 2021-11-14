# syntax=docker/dockerfile:1

 FROM golang:1.16-alpine

 WORKDIR /app

 COPY go.mod ./
 COPY go.sum ./
 RUN go mod download

 COPY *.go ./

 RUN go build -o /docker-go-service

 EXPOSE 8082

 CMD [ "/docker-go-service" ] 