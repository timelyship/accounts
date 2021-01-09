package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/url"
	"os"
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
		accountService.logger.Error("Sign up request validation error", zap.Error(validationError.Error))
		return validationError
	}
	// check if an user exists with the email
	if isExistingEmail, existingEmailError :=
		accountService.accountRepository.IsExistingEmail(signUpRequest.Email); existingEmailError != nil {
		accountService.logger.Error("isExistingEmail", zap.Error(existingEmailError.Error))
		return existingEmailError
	} else if isExistingEmail {
		bizError := fmt.Errorf("an user already exists with email %s", signUpRequest.Email)
		accountService.logger.Error("bizError", zap.Error(bizError))
		return utility.NewBadRequestError("Email Already exists", &bizError)
	}
	hashedPassword, passwordHashErr := utility.HashPassword(signUpRequest.Password)
	if passwordHashErr != nil {
		passwordHashErrorRest := fmt.Errorf("password hash error %v", passwordHashErr)
		accountService.logger.Error("passwordHashErrorRest", zap.Error(passwordHashErrorRest))
		return utility.NewBadRequestError("Email Already exists", &passwordHashErrorRest)
	}
	// create user
	newUserID := primitive.NewObjectID()
	user := domain.User{
		BaseEntity: domain.BaseEntity{
			ID: newUserID, InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		FirstName:       signUpRequest.FirstName,
		LastName:        signUpRequest.LastName,
		Email:           signUpRequest.Email,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		Active:          false,
		Password:        hashedPassword,
		Roles:           []*domain.Role{&domain.AppUserRole},
	}
	sErr := accountService.accountRepository.SaveUser(&user)
	if sErr != nil {
		accountService.logger.Info("User saved successfully", zap.Error(sErr.Error))
		return sErr
	}
	accountService.logger.Info("User saved successfully")
	emailVerErr := accountService.sendEmailVerificationMail(&user)
	if emailVerErr != nil {
		accountService.logger.Info("Email sending failed, inconsistent database state", zap.Error(emailVerErr.Error))
		return emailVerErr
	}
	accountService.logger.Info("Email sent successfully")
	accountService.createDefaultProfilePicture(newUserID.Hex())
	return nil
}

func (accountService *AccountService) sendEmailVerificationMail(user *domain.User) *utility.RestError {
	secret := utility.GetUUIDWithoutDash()
	vs := &domain.VerificationSecret{
		BaseEntity: domain.BaseEntity{
			ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		Type:       application.StringConst.Email,
		Subject:    user.Email,
		Secret:     secret,
		ValidUntil: time.Now().Add(time.Hour * 48),
	}
	err := accountService.accountRepository.SaveVerificationSecret(vs)
	if err != nil {
		accountService.logger.Info("Verification secret save failed", zap.Error(err.Error))
		return err
	}
	msgPayload := dto.NewEmailVerificationMsgPayload(
		application.StringConst.VerifyEmail, []string{user.Email},
		[]string{}, []string{"najim.ju@gmail.com"}, map[string]interface{}{
			"fullName":       fmt.Sprintf("%v %v", user.FirstName, user.LastName),
			"verificationId": secret,
		},
		[]string{}, "ahmedmdnajim@gmail.com",
	)
	bytes, mErr := json.Marshal(msgPayload)
	if mErr != nil {
		accountService.logger.Error("Unable to serialize email payload", zap.Error(mErr))
		return utility.NewInternalServerError("JSON serialization failed", &mErr)
	}
	pErr := utility.PublishEmailVerificationEvent(string(bytes))
	if pErr != nil {
		accountService.logger.Error("Unable to send payload to aws", zap.Error(pErr.Error))
		return pErr
	}
	return nil
}

func (accountService *AccountService) VerifyEmail(verificationToken string) *utility.RestError {
	verificationSecret, err := accountService.accountRepository.GetVerificationSecret(verificationToken)
	accountService.logger.Info("Verification secret",
		zap.String("verificationToken", verificationToken),
		zap.String("Subject", verificationSecret.Subject))
	if err != nil {
		accountService.logger.Error("Unable to fetch verification secret", zap.Error(err.Error))
		return err
	}
	user, dbErr := accountService.accountRepository.GetUserByEmail(verificationSecret.Subject)
	if dbErr != nil {
		accountService.logger.Error("Unable to fetch user by email ", zap.Error(dbErr.Error))
		return dbErr
	}
	user.Email = verificationSecret.Subject
	user.IsEmailVerified = true
	user.Active = true
	saveErr := accountService.accountRepository.UpdateUser(user)
	if saveErr != nil {
		accountService.logger.Error("Failed to update user after email verification", zap.Error(saveErr.Error))
	}
	// duplicate, to make sure user has a profile pic by default, increase probability
	accountService.createDefaultProfilePicture(user.ID.Hex())
	return saveErr
}

func (accountService *AccountService) createDefaultProfilePicture(userID string) {
	sess, sessErr := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)
	if sessErr != nil {
		accountService.logger.Error("Error creating default profile picture, aws session creation failed", zap.Error(sessErr))
		return
	}
	s3Svc := s3.New(sess)
	bucketName := os.Getenv("S3_BUCKET_PROFILE_PICTURE")
	src := fmt.Sprintf("/%s/%s", bucketName, "profile-default.png")
	copyObjectInput := &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(url.QueryEscape(src)),
		Key:        aws.String(userID),
	}
	result, err := s3Svc.CopyObject(copyObjectInput)
	if err != nil {
		accountService.logger.Error("Error creating default profile picture, copy failed", zap.Error(err))
	} else {
		accountService.logger.Info("Users default profile created", zap.Any("aws-s3-result", result))
	}

}
