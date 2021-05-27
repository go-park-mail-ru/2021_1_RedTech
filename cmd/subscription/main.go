package main

import (
	"Redioteka/internal/app/subscription"
	"Redioteka/internal/pkg/config"
)

func main() {
	subscription.RunServer(config.Get().Subscription.Port)
}
