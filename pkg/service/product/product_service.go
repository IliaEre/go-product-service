package product

import (
	d "aws-school-service/pkg/domain"
	"aws-school-service/pkg/repository"
	"io/ioutil"
	"log"

	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

type ProductService struct {
	repo *repository.DymanoRepository
}

func NewProductService(repository *repository.DymanoRepository) *ProductService {
	return &ProductService{repository}
}

func (s *ProductService) FindAll(w http.ResponseWriter, r *http.Request) {
	result, err := s.repo.FindAll()

	if err != nil {
		handleException(w, "problem...", http.StatusInternalServerError, err)
	}

	if len(result.Items) == 0 {
		handleException(w, "Not found", http.StatusNotFound, err)
	}

	products := []d.Product{}
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &products); err != nil {
		handleException(w, "Got error unmarshalling:", http.StatusInternalServerError, err)
	}

	resultEncoder(w, products)
}

func (hps *ProductService) FindOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("Try to find product with id: ", id)

	result, err := hps.repo.FindOne(id)
	if err != nil {
		handleException(w, "Query API call failed:", http.StatusInternalServerError, err)
	}

	var resultSet []d.Product
	for _, i := range result.Items {
		product := d.Product{}
		err = dynamodbattribute.UnmarshalMap(i, &product)

		if err != nil {
			handleException(w, "Got error unmarshalling:", http.StatusInternalServerError, err)
		}
		resultSet = append(resultSet, product)
	}

	resultEncoder(w, resultSet)
}

func (hps *ProductService) Create(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleException(w, "Problem with body!", http.StatusInternalServerError, err)
	}

	var product d.Product
	json.Unmarshal(requestBody, &product)

	rError := hps.repo.Create(product)

	log.Println("Product:", product)
	if rError != nil {
		handleException(w, "Got error calling PutItem:", http.StatusInternalServerError, err)
	}

	resultEncoder(w, product)
}

func handleException(w http.ResponseWriter, message string, code int, err error) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(code)
	log.Println(message, err)
	http.Error(w, message, code)
}

func resultEncoder(w http.ResponseWriter, obj interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}
