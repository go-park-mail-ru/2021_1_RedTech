package usecase

import (
	"Redioteka/internal/app/domain"
	"Redioteka/internal/app/user"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(u domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: u,
	}
}

func (uc *userUsecase) GetById(id uint) (domain.User, error) {
	return uc.userRepo.GetById(id)
}

func (uc *userUsecase) GetCurrent() (domain.User, error) {
	// access to session
	userId, err := getCurrentUserIdFromSession()
	if err != nil {
		return domain.User{}, user.UnauthorizedError
	}
	return uc.GetById(userId)
}

func (uc *userUsecase) Signup(u *domain.User) (uint, error) {
	id, err := uc.userRepo.Store(u)
	if err != nil {
		return 0, user.AlreadyAddedError
	}
	return id, nil
}

func (uc *userUsecase) Login(u *domain.User) (uint, error) {
	foundUser, err := uc.userRepo.GetByEmail(u.Email)
	if err != nil {
		return 0, user.NotFoundError
	}
	if foundUser.Password != u.Password {
		return 0, user.InvalidCredentials
	}
	return foundUser.ID, nil
}

func (uc *userUsecase) Logout(u *domain.User) error {
	return nil
}

func (uc *userUsecase) Update(updatedUser *domain.User) error {
	if !isUpdateValid(updatedUser) {
		return user.InvalidUpdateError
	}
	return uc.userRepo.Update(updatedUser)
}

func (uc *userUsecase) Delete(id uint) error {
	return uc.userRepo.Delete(id)
}

func isUpdateValid(update *domain.User) bool {
	return !(update.Email == "" && update.Username == "")
}

func getCurrentUserIdFromSession() (uint, error) {
	panic("implement me")
	return 1, nil
}
