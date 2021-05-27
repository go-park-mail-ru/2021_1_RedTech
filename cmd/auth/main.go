package main

import (
	"Redioteka/internal/app/auth"
	"Redioteka/internal/pkg/config"
)

func main() {
	auth.RunServer(config.Get().Auth.Port)
}
