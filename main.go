package main

import (
	"timelyship.com/accounts/config"
	"timelyship.com/accounts/repository"
)

func main() {
	defer repository.DisconnectMongoClient()
	config.Init()
	config.Start()
}
