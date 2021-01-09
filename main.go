package main

import (
	"timelyship.com/accounts/config"
)

func main() {
	config.Init()
	config.Start()
}
