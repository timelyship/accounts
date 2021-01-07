package response

import "timelyship.com/accounts/dto"

type ProfileResponse struct {
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Email           string      `json:"email"`
	Picture         string      `json:"picture"`
	Roles           []*dto.Role `json:"roles"`
	PhoneNumber     string      `json:"phone"`
	IsPhoneVerified bool        `json:"isPhoneVerified"`
	UserID          string      `json:"sub"`
}
