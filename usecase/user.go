package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"stock-control-back/model"
	"stock-control-back/repository"
	"stock-control-back/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	MongoCollection *mongo.Collection
}

func (service *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	response := &model.Response{}

	defer json.NewEncoder(w).Encode(response)

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid body")
		response.Error = err.Error()
		return
	}
	repo := repository.UserRepo{MongoCollection: service.MongoCollection}

	existingUser, err := repo.FindUserByEmail(user.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error finding user by email", err)
		response.Error = "Internal server error"
		return
	}

	if existingUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("User with the same email already exists")
		response.Error = "User with the same email already exists"
		return
	}

	user.ID = uuid.NewString()
	tokenString, err := utils.GenerateToken(user.ID, user.Email)

	CreateID, err := repo.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Insert error", err)
		response.Error = err.Error()
		return
	}

	response.Data = tokenString
	w.WriteHeader(http.StatusOK)

	log.Println("User created with id ", CreateID, user)
}
