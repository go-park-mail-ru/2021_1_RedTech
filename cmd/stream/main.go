package main

import (
	"Redioteka/internal/app/stream"
	"Redioteka/internal/pkg/config"
)

func main() {
	stream.RunServer(config.Get().Stream.Port)
}
