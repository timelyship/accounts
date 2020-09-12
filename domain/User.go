package domain

import (
	"regexp"
	"strings"
	"time"
	"timelyship.com/accounts/utility"
)

// https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type User struct {
	Id          int64
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	DateCreated time.Time
}

func Normalize(user *User) {
	user.Email = strings.TrimSpace(user.Email)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
}

// violation of database constraints are handled at the service layer, because those are business validations.
// Data validation goes here.
func ValidateUser(user *User) *utility.RestError {
	if !emailRegex.MatchString(user.Email) {
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
