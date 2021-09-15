package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Slogan struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Pelicula string             `json:"pelicula,omitempty" bson:"pelicula,omitempty"`
	Eslogan  string             `json:"eslogan,omitempty" bson:"eslogan,omitempty"`
	Contexto string             `json:"contexto,omitempty" bson:"contexto,omitempty"`
}
