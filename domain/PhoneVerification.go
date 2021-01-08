package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PhoneVerification struct {
	BaseEntity `bson:",inline"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Phone      string             `bson:"phone"`
}
