package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BaseEntity struct {
	Id         primitive.ObjectID `bson:"_id"`
	InsertedAt time.Time          `bson:"inserted_at"`
	LastUpdate time.Time          `bson:"last_update"`
}
