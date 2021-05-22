package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"crypto/sha256"
)

type authorizationUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthorizationUsecase(u domain.UserRepository) domain.AuthorizationUsecase {
	return &authorizationUsecase{
		userRepo: u,
	}
}

func preparePassword(u *domain.User) {
	u.Password = sha256.Sum256([]byte(u.InputPassword))
	u.InputPassword = ""
	u.ConfirmInputPassword = ""
}

func isSignupFormValid(uForm *domain.User) bool {
	return uForm.Username != "" && uForm.Email != "" && uForm.InputPassword != "" && uForm.InputPassword == uForm.ConfirmInputPassword
}

func (a authorizationUsecase) Signup(u *domain.User) (domain.User, error) {
	if !isSignupFormValid(u) {
		return domain.User{}, user.InvalidCredentials
	}

	preparePassword(u)
	id, err := a.userRepo.Store(u)
	if err != nil {
		return domain.User{}, user.AlreadyAddedError
	}

	createdUser, err := a.userRepo.GetById(id)
	if err != nil {
		return domain.User{}, user.NotFoundError
	}

	return createdUser, nil
}

func isLoginFormValid(uForm *domain.User) bool {
	return uForm.Email != "" && uForm.InputPassword != ""
}

func (a authorizationUsecase) Login(u *domain.User) (domain.User, error) {
	if !isLoginFormValid(u) {
		return domain.User{}, user.InvalidForm
	}

	foundUser, err := a.userRepo.GetByEmail(u.Email)
	if err != nil {
		return domain.User{}, user.NotFoundError
	}

	preparePassword(u)
	if foundUser.Password != u.Password {
		return domain.User{}, user.InvalidCredentials
	}

	return foundUser, nil
}

func (a authorizationUsecase) GetById(id uint) (domain.User, error) {
	return a.userRepo.GetById(id)
}

func (a authorizationUsecase) GetByEmail(email string) (domain.User, error) {
	return a.userRepo.GetByEmail(email)
}

func isUpdateValid(update *domain.User) bool {
	return update.Email != "" || update.Username != "" || update.Avatar != ""
}

func (a authorizationUsecase) Update(updatedUser *domain.User) error {
	if !isUpdateValid(updatedUser) {
		return user.InvalidUpdateError
	}
	return a.userRepo.Update(updatedUser)
}

func (a authorizationUsecase) Store(user *domain.User) (uint, error) {
	return a.userRepo.Store(user)
}

func (a authorizationUsecase) Delete(id uint) error {
	return a.userRepo.Delete(id)
}
