package request

import (
	"timelyship.com/accounts/utility"
)

type LoginRequest struct {
	EmailOrPhone string `json:"emailOrPhone"`
	Password     string `json:"password"`
}

func (r LoginRequest) ApplyUIValidation() *utility.RestError {
	if r.EmailOrPhone == "" {
		return utility.NewBadRequestError("EmailOrPhone field can not be empty", nil)
	}
	return nil
}
