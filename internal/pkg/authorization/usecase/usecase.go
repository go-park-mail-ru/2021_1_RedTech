package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authorizationUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthorizationUsecase(u domain.UserRepository) domain.AuthorizationUsecase {
	return &authorizationUsecase{
		userRepo: u,
	}
}

func (a authorizationUsecase) setSub(id uint, toSet *domain.User) {
	if time.Until(a.userRepo.CheckSub(id)) > 0 {
		toSet.IsSubscriber = true
	}
}

func (a authorizationUsecase) GetById(id uint) (domain.User, error) {
	log.Log.Info(fmt.Sprint("GetById", id))
	foundUser, err := a.userRepo.GetById(id)
	if err != nil {
		return domain.User{}, err
	}
	a.setSub(id, &foundUser)
	return foundUser, nil
}

func (a authorizationUsecase) GetByEmail(email string) (domain.User, error) {
	foundUser, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	a.setSub(foundUser.ID, &foundUser)
	return foundUser, nil
}

func preparePassword(u *domain.User) {
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(u.InputPassword), bcrypt.DefaultCost)
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

	if bcrypt.CompareHashAndPassword(foundUser.Password, []byte(u.InputPassword)) != nil {
		return domain.User{}, user.InvalidCredentials
	}
	return foundUser, nil
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

func (a authorizationUsecase) Delete(id uint) error {
	return a.userRepo.Delete(id)
}
