package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BaseEntity struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	InsertedAt time.Time          `json:"inserted_at" bson:"inserted_at"`
	LastUpdate time.Time          `json:"last_update" bson:"last_update"`
}

type Person struct {
	BaseEntity `bson:",inline"`
	Name       string    `json:"name" bson:"name"`
	BirthDate  time.Time `json:"birth_date" bson:"birth_date"`
}
