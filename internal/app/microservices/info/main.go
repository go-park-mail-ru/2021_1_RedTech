package main

import (
	server2 "Redioteka/internal/app/microservices/auth/app/server"
)

func main() {
	server2.RunServer(":8081")
}
