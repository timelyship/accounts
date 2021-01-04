package main

import (
	"context"
	"timelyship.com/accounts/config"
	"timelyship.com/accounts/repository"
)

func main() {
	config.Init()
	config.Start()
	repository.MongoClient.Disconnect(context.TODO())
}
