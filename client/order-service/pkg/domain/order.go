package domain

type Order struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Desc"`
	ProductID   string `json:"productId"`
}
