package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const USER_COLLECTION = "user"

//
//import (
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"go.mongodb.org/mongo-driver/mongo/readpref"
//	"time"
//	"timelyship.com/accounts/domain"
//	"timelyship.com/accounts/utility"
//)
//
//func SaveUser(user *domain.User) *utility.RestError {
//	uri := "mongodb+srv://mongodbroot:s3curedp%40s%24w0rd89@mongowork.sxfuk.mongodb.net/?retryWrites=true&w=majority"
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
//	if err != nil {
//		return utility.NewInternalServerError("Could not connect to database.Try after some time.")
//	}
//
//	defer func() {
//		if err = client.Disconnect(ctx); err != nil {
//			panic(err)
//		}
//	}()
//
//	// Ping the primary
//	if err := client.Ping(ctx, readpref.Primary()); err != nil {
//		panic(err)
//	}
//
//	fmt.Println("Successfully connected and pinged.")
//
//	insertResult, error := client.Database("timelyship-dev-db").Collection("users").InsertOne(ctx, user)
//	if error != nil {
//		fmt.Println("db-error:", error)
//		return utility.NewInternalServerError("Could not insert to database. Try after some time.")
//	}
//	fmt.Println("Successfully inserted", insertResult)
//	return nil
//}
//
//func GetById(id int64) (*domain.User, *utility.RestError) {
//	return nil, nil
//}
//
//func SavePerson(person *domain.Person) *utility.RestError {
//	uri := "mongodb+srv://mongodbroot:s3curedp%40s%24w0rd89@mongowork.sxfuk.mongodb.net/?retryWrites=true&w=majority"
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
//	if err != nil {
//		return utility.NewInternalServerError("Could not connect to database.Try after some time.")
//	}
//
//	defer func() {
//		if err = client.Disconnect(ctx); err != nil {
//			panic(err)
//		}
//	}()
//
//	// Ping the primary
//	if err := client.Ping(ctx, readpref.Primary()); err != nil {
//		panic(err)
//	}
//
//	fmt.Println("Successfully connected and pinged.")
//
//	insertResult, error := client.Database("timelyship-dev-db").Collection("users").InsertOne(ctx, person)
//	if error != nil {
//		fmt.Println("db-error:", error)
//		return utility.NewInternalServerError("Could not insert to database. Try after some time.")
//	}
//	fmt.Println("Successfully inserted", insertResult)
//	return nil
//}

func GetUserByGoogleId(googleId string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"google_auth_info.id", googleId}}
	result := domain.User{}
	error := GetCollection(USER_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not insert to database. Try after some time.")
	}
	return &result, nil
}

func SaveUser(user *domain.User) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(USER_COLLECTION).InsertOne(ctx, user)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.")
	}
	return nil
}
