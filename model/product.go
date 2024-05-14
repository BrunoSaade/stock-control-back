package model

type Product struct {
	ID       string  `json:"id,omitempty" bson:"id""`
	Name     string  `json:"name,omitempty" bson:"name""`
	Price    float64 `json:"price,omitempty" bson:"price""`
	Quantity int     `json:"quantity,omitempty" bson:"quantity""`
}
