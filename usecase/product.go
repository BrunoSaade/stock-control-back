package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"stock-control-back/model"
	"stock-control-back/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService struct {
	MongoCollection *mongo.Collection
}

func (service *ProductService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	response := &model.Response{}
	defer json.NewEncoder(w).Encode(response)

	var product model.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid body")
		response.Error = err.Error()
		return
	}

	userID := r.Context().Value("userID").(string)
	product.ID = uuid.NewString()

	repo := repository.ProductRepo{MongoCollection: service.MongoCollection}

	err_add := repo.AddProductToStock(userID, &product)
	if err_add != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Insert error", err)
		response.Error = err.Error()
		return
	}

	response.Data = product.ID
	w.WriteHeader(http.StatusOK)
}
func (service *ProductService) FindAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	response := &model.Response{}
	defer json.NewEncoder(w).Encode(response)

	repo := repository.ProductRepo{MongoCollection: service.MongoCollection}
	userID := r.Context().Value("userID").(string)

	products, err := repo.FindStockByUserID(userID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: ", err)
		response.Error = err.Error()
		return
	}

	response.Data = products
	w.WriteHeader(http.StatusOK)
}
func (service *ProductService) DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	response := &model.Response{}
	defer json.NewEncoder(w).Encode(response)

	productID := mux.Vars(r)["id"]
	log.Println("ProductID ", productID)

	repo := repository.ProductRepo{MongoCollection: service.MongoCollection}
	userID := r.Context().Value("userID").(string)

	err := repo.RemoveProductFromStock(userID, productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: ", err)
		response.Error = err.Error()
		return
	}

	response.Data = true
	w.WriteHeader(http.StatusOK)
}
