package grpc

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/subscription/delivery/grpc/proto"

	"golang.org/x/net/context"
)

type SubscriptionHandler struct {
	subUsecase domain.SubscriptionUsecase
}

func NewSubscriptionHandler(sub domain.SubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{
		subUsecase: sub,
	}
}

func (sh *SubscriptionHandler) Create(ctx context.Context, payment *proto.Payment) (*proto.ErrorMessage, error) {
	err := sh.subUsecase.Create(payment)
	if err != nil {
		return &proto.ErrorMessage{Error: 1}, err
	}
	return &proto.ErrorMessage{Error: 0}, nil
}

func (sh *SubscriptionHandler) Delete(ctx context.Context, user *proto.UserId) (*proto.ErrorMessage, error) {
	err := sh.subUsecase.Delete(uint(user.ID))
	if err != nil {
		return &proto.ErrorMessage{Error: 1}, err
	}
	return &proto.ErrorMessage{Error: 0}, nil
}
