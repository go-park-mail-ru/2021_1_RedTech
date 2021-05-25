package main

import (
	"Redioteka/internal/app/auth"
	"Redioteka/internal/constants"
)

func main() {
	auth.RunServer(constants.AuthServiceAddress)
}
