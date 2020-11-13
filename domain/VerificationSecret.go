package domain

import "time"

type VerificationSecret struct {
	BaseEntity `bson:",inline"`
	Type       string `json:"type" bson:"type"`
	// could be email or phone number
	Subject    string    `json:"subject" bson:"subject"`
	Secret     string    `json:"secret" bson:"secret"`
	ValidUntil time.Time `json:"validUntil" bson:"valid_until"`
}
