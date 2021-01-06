package domain

import (
	"time"
)

// https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
// var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?
// (?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

/*
type PhoneNumber struct {
	Number      string `bson:"number"`
	IsPrimary   bool   `bson:"is_primary"`
	IsVerified  bool   `bson:"is_verified"`
	IsEmergency bool   `bson:"is_verified"`
}
*/

type GoogleAuthInfo struct {
	ID      string
	Email   string // email must be verified
	Picture string
}

type FacebookAuthInfo struct {
	ID      string
	Email   string // email must be verified
	Picture string
}

type Role struct {
	Name string `bson:"name"`
}

var AppUserRole = Role{
	Name: "APP_USER",
}

type User struct {
	BaseEntity       `bson:",inline"`
	FirstName        string           `bson:"first_name"`
	LastName         string           `bson:"last_name"`
	Email            string           `bson:"email"`
	IsEmailVerified  bool             `bson:"is_email_verified"`
	IsPhoneVerified  bool             `bson:"is_phone_verified"`
	PhoneNumber      string           `bson:"phone"`
	Picture          string           `bson:"picture"`
	DateCreated      time.Time        `bson:"date_created"`
	DateUpdated      time.Time        `bson:"date_updated"`
	GoogleAuthInfo   GoogleAuthInfo   `bson:"google_auth_info"`
	FacebookAuthInfo FacebookAuthInfo `bson:"facebook_auth_info"`
	Password         string           `bson:"password"`
	Roles            []*Role          `bson:"roles"`
	Active           bool             `bson:"is_active"`
}
