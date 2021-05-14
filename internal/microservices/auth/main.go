package main

import (
	"Redioteka/internal/microservices/auth/app/server"
)

func main() {
	server.RunServer(":8082")
}
