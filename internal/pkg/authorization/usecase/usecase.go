package usecase

import "Redioteka/internal/pkg/domain"

type authorizationUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthorizationUsecase(u domain.UserRepository) domain.AuthorizationUsecase {
	return &authorizationUsecase{
		userRepo: u,
	}
}

func (a authorizationUsecase) GetById(id uint) (domain.User, error) {
	return a.userRepo.GetById(id)
}

func (a authorizationUsecase) GetByEmail(email string) (domain.User, error) {
	return a.userRepo.GetByEmail(email)
}

func (a authorizationUsecase) Update(user *domain.User) error {
	panic("implement me")
}

func (a authorizationUsecase) Store(user *domain.User) (uint, error) {
	return a.userRepo.Store(user)
}

func (a authorizationUsecase) Delete(id uint) error {
	return a.userRepo.Delete(id)
}
