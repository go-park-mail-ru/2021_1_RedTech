package usecase

import "Redioteka/internal/pkg/domain"

type streamUsecase struct {
	streamRepo domain.StreamRepository
}

func NewStreamUsecase(s domain.StreamRepository) domain.StreamUsecase {
	return &streamUsecase{
		streamRepo: s,
	}
}

func (s streamUsecase) GetStream(id uint) ([]domain.Stream, error) {
	return s.streamRepo.GetStream(id)
}
