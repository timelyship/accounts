package domain

type PhoneVerification struct {
	BaseEntity `bson:",inline"`
	UserID     string `bson:"user_id"`
	Phone      string `bson:"phone"`
}
