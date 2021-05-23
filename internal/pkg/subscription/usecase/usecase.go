package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/subscription/delivery/grpc/proto"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/payment"
)

type subscriptionUsecase struct {
	subRepo domain.SubscriptionRepository
}

func NewSubscriptionUsecase(repo domain.SubscriptionRepository) *subscriptionUsecase {
	return &subscriptionUsecase{
		subRepo: repo,
	}
}

func (su *subscriptionUsecase) Create(form *proto.Payment) error {
	err := payment.Check(form)
	if err != nil {
		log.Log.Error(err)
		return err
	}

	sub := &domain.Subscription{}
	return su.subRepo.Create(sub)
}

func (su *subscriptionUsecase) Delete(id uint) error {
	return su.subRepo.Delete(id)
}
