package repository

import (
	"context"
	"stock-control-back/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	MongoCollection *mongo.Collection
}

func (r *UserRepo) CreateUser(user *model.User) (interface{}, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	result, err := r.MongoCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	_, err = r.MongoCollection.UpdateOne(context.Background(),
		bson.M{"id": user.ID},
		bson.M{"$set": bson.M{"stock": []model.Product{}}})
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *UserRepo) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.MongoCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
