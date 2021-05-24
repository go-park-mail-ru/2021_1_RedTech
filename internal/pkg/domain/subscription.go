package domain

import (
	"Redioteka/internal/pkg/subscription/delivery/grpc/proto"
	"time"
)

const Cost = 100

type Subscription struct {
	UserID    uint
	Expiraton time.Time
	Actual    bool
}

//go:generate mockgen -destination=../subscription/repository/mock/mock_repo.go -package=mock Redioteka/internal/pkg/domain SubscriptionRepository
type SubscriptionRepository interface {
	Create(sub *Subscription) error
	Delete(id uint) error
}

//go:generate mockgen -destination=../subscription/usecase/mock/mock_usecase.go -package=mock Redioteka/internal/pkg/domain SubscriptionUsecase
type SubscriptionUsecase interface {
	Create(form *proto.Payment) error
	Delete(id uint) error
}
