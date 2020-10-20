package domain

import (
	"regexp"
	"time"
)

// https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type PhoneNumber struct {
	Number     string
	IsPrimary  bool
	IsVerified bool
}

type GoogleAuthInfo struct {
	Id      string
	Email   string // email must be verified
	Picture string
}

type Role struct {
	Name   string
	Parent *Role
}

type User struct {
	BaseEntity             `bson:",inline"`
	FirstName              string         `json:"firstName" bson:"first_name"`
	LastName               string         `json:"lastName" bson:"last_name"`
	PrimaryEmail           string         `json:"primaryEmail" bson:"primary_email"`
	IsPrimaryEmailVerified bool           `json:"isPrimaryEmailVerified" bson:"is_primary_email_verified"`
	PrimaryPicture         string         `json:"PrimaryPicture" bson:"primary_picture"`
	PhoneNumbers           []PhoneNumber  `json:"phoneNumbers" bson:"phone_numbers"`
	DateCreated            time.Time      `json:"dateCreated" bson:"date_created"`
	DateUpdated            time.Time      `json:"dateUpdated" bson:"date_updated"`
	GoogleAuthInfo         GoogleAuthInfo `json:"googleAuthInfo" bson:"google_auth_info"`
}

/*
func Normalize(user *User) {
	user.PrimaryEmail = strings.TrimSpace(user.PrimaryEmail)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
}

// violation of database constraints are handled at the service layer, because those are business validations.
// Data validation goes here.
func ValidateUser(user *User) *utility.RestError {
	if !emailRegex.MatchString(user.PrimaryEmail) {
		return utility.NewBadRequestError("Invalid email pattern. Please contact administrator.")
	}
	firstNameLen := len(user.FirstName)
	if firstNameLen == 0 {
		return utility.NewBadRequestError("First name can not be empty.")
	}
	if firstNameLen > 16 {
		return utility.NewBadRequestError("First name length can not exceed 16 characters")
	}
	lastNameLen := len(user.LastName)
	if lastNameLen > 16 {
		return utility.NewBadRequestError("Last name length can not exceed 16 characters")
	}
	return nil
}

*/
