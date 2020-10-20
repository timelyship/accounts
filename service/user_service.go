package service

//
//import (
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"time"
//	"timelyship.com/accounts/domain"
//	user2 "timelyship.com/accounts/repository/user"
//	"timelyship.com/accounts/utility"
//)
//
//func GetUser() {
//
//}
//
//func CreateUser(user domain.User) (*domain.User, *utility.RestError) {
//	if user.Id == 100 {
//		return nil, utility.NewInternalServerError("Failed to save user.")
//	}
//	user2.SaveUser(&user)
//	person := domain.Person{
//		BaseEntity: domain.BaseEntity{Id: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
//		Name:       "John Belushi",
//		BirthDate:  time.Date(1959, time.February, 28, 0, 0, 0, 0, time.UTC),
//	}
//	user2.SavePerson(&person)
//	return &user, nil
//}
//
//func FindUser() {
//
//}
