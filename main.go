package main

import (
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/config"
)

func main() {
	config.Init()
	application.Start()
}
