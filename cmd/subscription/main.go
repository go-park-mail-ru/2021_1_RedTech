package main

import (
	"Redioteka/internal/app/subscription"
	"Redioteka/internal/constants"
)

func main() {
	subscription.RunServer(constants.SubscriptionServiceAddress)
}
