package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

type ProfileService struct {
	logger            zap.Logger
	profileRepository repository.ProfileRepository
}

var (
	fields = map[string]string{
		"firstName":      "first_name",
		"lastName":       "last_name",
		"phoneNums":      "phone_numbers",
		"profilePicture": "primary_picture",
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

func ProvideProfileService(l zap.Logger, p repository.ProfileRepository) ProfileService {
	return ProfileService{
		logger:            l,
		profileRepository: p,
	}
}
