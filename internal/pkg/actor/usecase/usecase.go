package usecase

import "Redioteka/internal/pkg/domain"

type actorUsecase struct {
	actorRepo domain.ActorRepository
}

func NewActorUsecase(ar domain.ActorRepository) domain.ActorUsecase {
	return &actorUsecase{
		actorRepo: ar,
	}
}

func (a actorUsecase) GetById(id uint) (domain.Actor, error) {
	return a.actorRepo.GetById(id)
}
