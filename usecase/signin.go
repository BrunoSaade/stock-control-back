package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"stock-control-back/model"
	"stock-control-back/repository"
	"stock-control-back/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SigninService struct {
	MongoCollection *mongo.Collection
}

func (service *SigninService) Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	response := &model.Response{}

	defer json.NewEncoder(w).Encode(response)

	var signinRequest model.SigninRequest

	err := json.NewDecoder(r.Body).Decode(&signinRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid body")
		response.Error = err.Error()
		return
	}

	userRepo := repository.UserRepo{MongoCollection: service.MongoCollection}
	user, err := userRepo.FindUserByEmail(signinRequest.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error finding user:", err)
		response.Error = "Internal server error"
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("User not found")
		response.Error = "Invalid username or password"
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signinRequest.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Invalid password")
		response.Error = "Invalid username or password"
		return
	}

	tokenString, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to generate token")
		response.Error = "Failed to generate token"
		return
	}

	response.Data = tokenString
	w.WriteHeader(http.StatusOK)

	log.Println("User logged in successfully")
}
