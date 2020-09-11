package domain

import (
	"time"
)

type User struct {
	Id          int64
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	DateCreated time.Time
}
