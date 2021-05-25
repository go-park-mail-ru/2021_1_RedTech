package main

import (
	"Redioteka/internal/app/server"
	"Redioteka/internal/pkg/config"
)

func main() {
	server.RunServer(config.Get().Main.Port)
}
