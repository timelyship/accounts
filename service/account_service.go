package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strings"
	"time"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

type AccountService struct {
	accountRepository repository.AccountRepository
	logger            zap.Logger
}

func ProvideAccountService(a repository.AccountRepository, z zap.Logger) AccountService {
	return AccountService{
		accountRepository: a,
		logger:            z,
	}
}

func (accountService *AccountService) InitiateSignUp(signUpRequest request.SignUpRequest) *utility.RestError {
	validationError := signUpRequest.ApplyUIValidation()
	if validationError != nil {
		accountService.logger.Error("Sign up request validation error", zap.Any("validation error", validationError))
		return validationError
	}
	// check if an user exists with the email
	if isExistingEmail, error := repository.IsExistingEmail(signUpRequest.Email); error != nil {
		accountService.logger.Error("isExistingEmail", zap.Any("isExistingEmail error", error))
		return error
	} else if isExistingEmail {
		bizError := fmt.Errorf("an user already exists with email %s", signUpRequest.Email)
		accountService.logger.Error("bizError", zap.Any("Email already exists", bizError.Error()))
		return utility.NewBadRequestError("Email Already exists", &bizError)
	}
	hashedPassword, passwordHashErr := utility.HashPassword(signUpRequest.Password)
	if passwordHashErr != nil {
		passwordHashErrorRest := fmt.Errorf("password hash error %v", passwordHashErr)
		accountService.logger.Error("passwordHashErrorRest", zap.Any("Email already exists", passwordHashErrorRest.Error()))
		return utility.NewBadRequestError("Email Already exists", &passwordHashErrorRest)
	}
	// create user
	user := domain.User{
		BaseEntity: domain.BaseEntity{
			ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		FirstName:              signUpRequest.FirstName,
		LastName:               signUpRequest.LastName,
		PrimaryEmail:           signUpRequest.Email,
		IsPrimaryEmailVerified: false,
		Password:               hashedPassword,
		Roles:                  []*domain.Role{&domain.AppUserRole},
	}
	sErr := repository.SaveUser(&user)
	if sErr != nil {
		return sErr
	}
	accountService.logger.Info("User saved successfully")
	emailVerErr := sendEmailVerificationMail(&user)
	if emailVerErr != nil {
		fmt.Println("Inconsistent DB Error")
		return emailVerErr
	}
	return nil
}

func sendEmailVerificationMail(user *domain.User) *utility.RestError {
	secret := strings.ReplaceAll(uuid.New().String(), "-", "")
	vs := &domain.VerificationSecret{
		BaseEntity: domain.BaseEntity{
			ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		Type:       application.StringConst.Email,
		Subject:    user.PrimaryEmail,
		Secret:     secret,
		ValidUntil: time.Now().Add(time.Hour * 48),
	}
	err := repository.SaveVerificationSecret(vs)
	if err != nil {
		return err
	}
	msgPayload := dto.NewEmailVerificationMsgPayload(
		application.StringConst.VerifyEmail, []string{user.PrimaryEmail},
		[]string{}, []string{"najim.ju@gmail.com"}, map[string]interface{}{
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

func VerifyEmail(verificationToken string) *utility.RestError {
	verificationSecret, err := repository.GetVerificationSecret(verificationToken)
	if err != nil {
		return err
	}
	user, dbErr := repository.GetUserByEmail(verificationSecret.Subject)
	if dbErr != nil {
		return dbErr
	}
	user.PrimaryEmail = verificationSecret.Subject
	user.IsPrimaryEmailVerified = true
	saveErr := repository.UpdateUser(user)
	return saveErr
}
