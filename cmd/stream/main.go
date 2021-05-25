package main

import (
	"Redioteka/internal/app/stream"
	"Redioteka/internal/constants"
)

func main() {
	stream.RunServer(constants.StreamServiceAddress)
}
