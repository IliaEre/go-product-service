package main

type Product struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type ResponseEntity struct {
	Product Product `json:"product"`
	TraceId string  `json:"traceId"`
}

type ResponseEntityAll struct {
	Products []Product `json:"products"`
	TraceId  string    `json:"traceId"`
}

var Products []Product
