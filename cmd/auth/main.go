package main

import (
	"Redioteka/internal/app/info"
	"Redioteka/internal/pkg/config"
)

func main() {
	info.RunServer(config.Get().Auth.Port)
}
