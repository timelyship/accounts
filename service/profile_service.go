package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
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
	// todo , should be a utility function
	sess, sessErr := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)
	if sessErr != nil {
		return nil, utility.NewInternalServerError("Could not create new aws session", &sessErr)
	}
	s3Svc := s3.New(sess)
	bucketName := os.Getenv("S3_BUCKET_PROFILE_PICTURE")
	url, urlErr := s.generateSignedUrl(s3Svc, bucketName, userID)
	if urlErr != nil {
		msg := fmt.Sprintf("Failed to generate pre signed url for user %v", userID)
		s.logger.Error(msg, zap.Error(*urlErr))
		url = ""
	}

	return &response.ProfileResponse{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Picture:         url,
		Roles:           convertRoles(user.Roles),
		PhoneNumber:     user.PhoneNumber,
		IsPhoneVerified: user.IsPhoneVerified,
		UserID:          userID,
	}, nil
}

func (s *ProfileService) ChangePhoneNumber(id, phone string) *utility.RestError {
	userID, parseHexErr := primitive.ObjectIDFromHex(id)
	if parseHexErr != nil {
		s.logger.Error("User id parse error", zap.Error(parseHexErr))
		return utility.NewInternalServerError("Could not parse userId", &parseHexErr)
	}

	err := s.profileRepository.ChangePhoneNumber(userID, phone)
	if err != nil {
		return err
	}
	phoneVerification := &domain.PhoneVerification{
		BaseEntity: domain.BaseEntity{
			ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		UserID: userID,
		Phone:  phone,
	}
	qErr := s.profileRepository.EnqueuePhoneVerification(userID, phoneVerification)
	if qErr != nil {
		s.logger.Warn("Inconsistent database state", zap.Error(qErr.Error))
	}
	return qErr
}

func (s *ProfileService) UploadProfilePhoto(id string, header *multipart.FileHeader, buf *bytes.Buffer) (*response.PhotoUploadResponse, *utility.RestError) {
	s.logger.Info("Uploading file to s3",
		zap.Int64("Size", header.Size), zap.String("Filename", header.Filename), zap.Any("Header", header.Header),
	)
	//extension, extErr := s.getFileExtension(header.Filename)
	//if extErr != nil {
	//	return nil, utility.NewBadRequestError("Bad file name", &extErr)
	//}
	sess, sessErr := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)
	if sessErr != nil {
		return nil, utility.NewInternalServerError("Could not create new aws session", &sessErr)
	}

	s3Svc := s3.New(sess)
	bucketName := os.Getenv("S3_BUCKET_PROFILE_PICTURE")
	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &id,
		Body:   buf,
	}
	uploader := s3manager.NewUploaderWithClient(s3Svc)
	result, err := uploader.Upload(upParams)
	s.logger.Info("result", zap.Any("result", result))
	if err != nil {
		return nil, utility.NewInternalServerError("s3 upload failed", &err)
	}
	response := &response.PhotoUploadResponse{}
	url, urlErr := s.generateSignedUrl(s3Svc, bucketName, id)
	if urlErr != nil {
		s.logger.Error("Fatal error, image uploaded to s3", zap.Error(*urlErr))
		response.Url = ""
	} else {
		s.logger.Debug("Signed Url success", zap.String("key", id))
		response.Url = url
	}
	//s.updateProfilePicture(id, key)
	return response, nil
}

func (s *ProfileService) generateSignedUrl(svc *s3.S3, bucketName string, key string) (string, *error) {
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign((7 * 24) * time.Hour) // 7 days,max from aws
	if err != nil {
		return "", &err
	}
	return urlStr, nil
}

//func (s *ProfileService) getFileExtension(filename string) (string, error) {
//	idx := strings.LastIndex(filename, ".")
//	if idx == -1 || idx == len(filename)-1 {
//		return "", errors.New(fmt.Sprintf("failed to discover file extension, last index of dot(.) is %v", idx))
//	}
//	return filename[idx+1:], nil
//}

//func (s *ProfileService) updateProfilePicture(id string, key string) {
//	userID, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		s.logger.Error("Sever error,failed to update")
//	}
//	s.profileRepository.UpdateImageUrl()
//}

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
