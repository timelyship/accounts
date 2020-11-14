package request

import (
	"timelyship.com/accounts/utility"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (r LoginRequest) ApplyUiValidation() *utility.RestError {
	if r.Email == "" && r.Phone == "" {
		return utility.NewBadRequestError("Both email and phone number can not be empty", nil)
	}
	return nil
}
