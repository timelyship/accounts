package domain

import "time"

type VerificationSecret struct {
	BaseEntity `bson:",inline"`
	Type       string    `bson:"type"`
	Subject    string    `bson:"subject"` // could be email or phone number
	Secret     string    `bson:"secret"`
	ValidUntil time.Time `bson:"valid_until"`
}
