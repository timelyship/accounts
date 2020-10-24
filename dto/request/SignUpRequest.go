package request

import (
	"errors"
	"fmt"
	"regexp"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/utility"
)

type SignUpRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (r SignUpRequest) Validate() *utility.RestError {
	validationErrors := make([]error, 0)
	/* Validate first name*/
	if len(r.FirstName) == 0 {
		validationErrors = append(validationErrors, errors.New("firstName should not be empty"))
	}
	if len(r.LastName) > application.INT_CONST.FIRST_NAME_MAX_LEN {
		validationErrors = append(validationErrors, errors.New(
			fmt.Sprintf("firstName should be within[0,%v] characters", application.INT_CONST.FIRST_NAME_MAX_LEN)))
	}
	/* Validate last name*/
	if len(r.LastName) == 0 {
		validationErrors = append(validationErrors, errors.New("lastName should not be empty"))
	}
	if len(r.LastName) > application.INT_CONST.LAST_NAME_MAX_LEN {
		validationErrors = append(validationErrors, errors.New(
			fmt.Sprintf("lastName should be within[0,%v] characters", application.INT_CONST.LAST_NAME_MAX_LEN)))
	}
	/* Validate email*/
	if len(r.Email) == 0 {
		validationErrors = append(validationErrors, errors.New("email should not be empty"))
	}
	if len(r.Email) > application.INT_CONST.EMAIL_NAME_MAX_LEN {
		validationErrors = append(validationErrors, errors.New(
			fmt.Sprintf("lastName should be within[0,%v] characters", application.INT_CONST.EMAIL_NAME_MAX_LEN)))
	}
	emailRegex := regexp.MustCompile(application.STRING_CONST.EMAIL_PATTERN)
	if !emailRegex.MatchString(r.Email) {
		validationErrors = append(validationErrors, errors.New("email pattern invalid,contact administrator"))
	}
	var err error = utility.ValidationError{
		ErrorMessages: validationErrors,
	}
	if len(validationErrors) > 0 {
		return utility.NewBadRequestError("Request body is invalid", &err)
	}
	return nil
}
