package repository

import (
	"context"
	"stock-control-back/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	MongoCollection *mongo.Collection
}

func (r *ProductRepo) AddProductToStock(userID string, product *model.Product) error {
	_, err := r.MongoCollection.UpdateOne(context.Background(),
		bson.M{"id": userID},
		bson.M{"$push": bson.M{"stock": product}})
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) FindStockByUserID(userID string) ([]model.Product, error) {
	var user model.User
	err := r.MongoCollection.FindOne(context.Background(), bson.M{"id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Stock, nil
}

func (r *ProductRepo) RemoveProductFromStock(userID string, productID string) error {
	_, err := r.MongoCollection.UpdateOne(context.Background(),
		bson.M{"id": userID},
		bson.M{"$pull": bson.M{"stock": bson.M{"id": productID}}})
	if err != nil {
		return err
	}
	return nil
}

// func (r *ProductRepo) DeleteProductByID(productID string) (int64, error) {
// 	result, err := r.MongoCollection.DeleteOne(context.Background(),
// 		bson.D{{Key: "id", Value: productID}})

// 	if err != nil {
// 		return 0, err
// 	}

// 	return result.DeletedCount, nil
// }
