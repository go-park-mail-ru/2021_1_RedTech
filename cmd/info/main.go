package main

import (
	"Redioteka/internal/app/info"
	"Redioteka/internal/constants"
)

func main() {
	info.RunServer(constants.InfoServiceAddress)
}
