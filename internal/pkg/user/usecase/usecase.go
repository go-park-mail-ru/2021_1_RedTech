package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/session"
	"io"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo   domain.UserRepository
	avatarRepo domain.AvatarRepository
}

func NewUserUsecase(u domain.UserRepository, a domain.AvatarRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo:   u,
		avatarRepo: a,
	}
}

func (uc *userUsecase) GetById(id uint) (domain.User, error) {
	user, err := uc.userRepo.GetById(id)
	if err != nil {
		return domain.User{}, err
	}
	if uc.userRepo.CheckSub(id).Sub(time.Now()) > 0 {
		user.IsSubscriber = true
	}
	return user, nil
}

func preparePassword(u *domain.User) {
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(u.InputPassword), bcrypt.DefaultCost)
	u.InputPassword = ""
	u.ConfirmInputPassword = ""
}

func isSignupFormValid(uForm *domain.User) bool {
	return uForm.Username != "" && uForm.Email != "" && uForm.InputPassword != "" && uForm.InputPassword == uForm.ConfirmInputPassword
}

func (uc *userUsecase) Signup(u *domain.User) (domain.User, error) {
	if !isSignupFormValid(u) {
		return domain.User{}, user.InvalidCredentials
	}

	preparePassword(u)
	id, err := uc.userRepo.Store(u)
	if err != nil {
		return domain.User{}, user.AlreadyAddedError
	}

	createdUser, err := uc.userRepo.GetById(id)
	if err != nil {
		return domain.User{}, user.NotFoundError
	}

	return createdUser, nil
}

func isLoginFormValid(uForm *domain.User) bool {
	return uForm.Email != "" && uForm.InputPassword != ""
}

func (uc *userUsecase) Login(u *domain.User) (domain.User, error) {
	if !isLoginFormValid(u) {
		return domain.User{}, user.InvalidForm
	}

	foundUser, err := uc.userRepo.GetByEmail(u.Email)
	if err != nil {
		return domain.User{}, user.NotFoundError
	}

	if bcrypt.CompareHashAndPassword(foundUser.Password, []byte(u.InputPassword)) != nil {
		return domain.User{}, user.InvalidCredentials
	}

	return foundUser, nil
}

func (uc *userUsecase) Logout(sess *session.Session) error {
	err := session.Manager.Delete(sess)
	if err != nil {
		return user.UnauthorizedError
	}
	return nil
}

func isUpdateValid(update *domain.User) bool {
	return update.Email != "" || update.Username != "" || update.Avatar != ""
}

func (uc *userUsecase) Update(updatedUser *domain.User) error {
	if !isUpdateValid(updatedUser) {
		return user.InvalidUpdateError
	}
	return uc.userRepo.Update(updatedUser)
}

func (uc *userUsecase) Delete(id uint) error {
	err := uc.userRepo.Delete(id)
	if err != nil {
		err = user.NotFoundError
	}
	return err
}

func (uc *userUsecase) GetFavourites(id uint, sess *session.Session) ([]domain.Movie, error) {
	err := session.Manager.Check(sess)
	if err != nil {
		return nil, user.UnauthorizedError
	}
	if sess.UserID != id {
		return nil, user.InvalidCredentials
	}

	return uc.userRepo.GetFavouritesByID(id)
}

func (uc *userUsecase) UploadAvatar(reader io.Reader, path, ext string) (string, error) {
	return uc.avatarRepo.UploadAvatar(reader, path, ext)
}
