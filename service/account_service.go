package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

func InitiateSignUp(signUpRequest request.SignUpRequest) *utility.RestError {
	validationError := signUpRequest.ApplyUiValidation()
	if validationError != nil {
		return validationError
	}
	// check if an user exists with the email
	if isExistingEmail, error := repository.IsExistingEmail(signUpRequest.Email); error != nil {
		return error
	} else if isExistingEmail {
		bizError := errors.New(fmt.Sprintf("An user already exists with email %s", signUpRequest.Email))
		return utility.NewBadRequestError("Email Already exists", &bizError)
	}
	// create user
	user := domain.User{
		FirstName:              signUpRequest.FirstName,
		LastName:               signUpRequest.LastName,
		PrimaryEmail:           signUpRequest.Email,
		IsPrimaryEmailVerified: false,
		Password:               utility.HashPassword(signUpRequest.Password),
	}
	repository.SaveUser(&user)
	emailVerErr := sendEmailVerificationMail(&user)
	if emailVerErr != nil {
		fmt.Println("Inconsistent DB Error")
		return emailVerErr
	}
	return nil
}

func sendEmailVerificationMail(user *domain.User) *utility.RestError {
	secret := uuid.New().String()
	vs := &domain.VerificationSecret{
		Type:       application.STRING_CONST.EMAIL,
		Subject:    user.PrimaryEmail,
		Secret:     secret,
		ValidUntil: time.Now().Add(time.Hour * 48),
	}
	err := repository.SaveVerificationSecret(vs)
	if err != nil {
		return err
	}
	msgPayload := dto.NewEmailVerificationMsgPayload(
		application.STRING_CONST.VERIFY_EMAIL, []string{user.PrimaryEmail}, []string{}, []string{"najim.ju@gmail.com"}, map[string]interface{}{
			"fullName":       fmt.Sprintf("%v %v", user.FirstName, user.LastName),
			"verificationId": secret,
		},
		[]string{}, "ahmedmdnajim@gmail.com",
	)
	bytes, mErr := json.Marshal(msgPayload)
	if mErr != nil {
		return utility.NewInternalServerError("JSON serialization failed", &mErr)
	}
	pErr := utility.PublishEmailVerificationEvent(string(bytes))
	if pErr != nil {
		return pErr
	}
	return nil
}
