package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(u domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: u,
	}
}

func (u *userUsecase) GetById(id uint) (domain.User, error) {
	return u.userRepo.GetById(id)
}

func (u *userUsecase) GetByEmail(email string) (domain.User, error) {
	return u.userRepo.GetByEmail(email)
}

func isUpdateValid(update *domain.User) bool {
	return !(update.Email == "" && update.Username == "")
}

func (u *userUsecase) Update(updatedUser *domain.User) error {
	if !isUpdateValid(updatedUser) {
		return user.InvalidUpdateError
	}
	return u.userRepo.Update(updatedUser)
}

func (u *userUsecase) Delete(id uint) error {
	return u.userRepo.Delete(id)
}
