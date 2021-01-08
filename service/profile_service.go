package service

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/dto/response"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

type ProfileService struct {
	logger            zap.Logger
	profileRepository repository.ProfileRepository
}

var (
	fields = map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
		"phone":     "phone",
		"picture":   "picture",
	}
)

// can update firstName,lastName, profilePicture and PhoneNumber
func (s *ProfileService) Patch(userID string, request []*request.ProfilePatchRequest) *utility.RestError {
	for _, r := range request {
		dbField, ok := fields[r.Field]
		if !ok {
			em := fmt.Sprintf("Invalid field for patch %s", r.Field)
			valErr := errors.New(em)
			s.logger.Error("validation error", zap.Error(valErr))
			return utility.NewBadRequestError(em, &valErr)
		}
		r.Field = dbField
	}
	return s.profileRepository.Patch(userID, request)
}

func (s *ProfileService) GetProfileById(userID string) (*response.ProfileResponse, *utility.RestError) {
	u, error := primitive.ObjectIDFromHex(userID)
	if error != nil {
		s.logger.Error("Error decoding hex", zap.Error(error))
		return nil, utility.NewInternalServerError(fmt.Sprintf("Can not convert %v  primitive.ObjectIDFromHex", userID), &error)
	}
	user, err := s.profileRepository.GetProfileById(u)
	if err != nil {
		s.logger.Error("Error fetching user by id", zap.Error(err.Error))
		return nil, err
	}
	return &response.ProfileResponse{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Picture:         user.Picture,
		Roles:           convertRoles(user.Roles),
		PhoneNumber:     user.PhoneNumber,
		IsPhoneVerified: user.IsPhoneVerified,
		UserID:          userID,
	}, nil
}

func (s *ProfileService) ChangePhoneNumber(id, phone string) *utility.RestError {
	err := s.profileRepository.ChangePhoneNumber(id, phone)
	if err != nil {
		return err
	}
	phoneVerification := &domain.PhoneVerification{
		BaseEntity: domain.BaseEntity{
			ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		UserID: id,
		Phone:  phone,
	}
	qErr := s.profileRepository.EnqueuePhoneVerification(phoneVerification)
	if qErr != nil {
		s.logger.Warn("Inconsistent database state", zap.Error(qErr.Error))
	}
	return qErr
}

func convertRoles(roles []*domain.Role) []*dto.Role {
	dtoRoles := make([]*dto.Role, 0)
	for _, r := range roles {
		dtoRoles = append(dtoRoles, &dto.Role{
			Name: r.Name,
		})
	}
	return dtoRoles
}

func ProvideProfileService(l zap.Logger, p repository.ProfileRepository) ProfileService {
	return ProfileService{
		logger:            l,
		profileRepository: p,
	}
}
