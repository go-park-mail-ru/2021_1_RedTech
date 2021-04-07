package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/session"
	"crypto/sha256"
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

func preparePassword(u *domain.User) {
	u.Password = sha256.Sum256([]byte(u.InputPassword))
	u.InputPassword = ""
	u.ConfirmInputPassword = ""
}

func isSignupFormValid(uForm *domain.User) bool {
	return uForm.Username != "" && uForm.Email != "" && uForm.InputPassword != "" && uForm.InputPassword == uForm.ConfirmInputPassword
}

func (uc *userUsecase) Signup(u *domain.User) (domain.User, *session.Session, error) {
	if !isSignupFormValid(u) {
		return domain.User{}, nil, user.InvalidCredentials
	}

	preparePassword(u)
	id, err := uc.userRepo.Store(u)
	if err != nil {
		return domain.User{}, nil, user.AlreadyAddedError
	}

	sess := &session.Session{UserID: id}
	err = session.Manager.Create(sess)
	if err != nil {
		return domain.User{}, nil, err
	}

	createdUser, err := uc.userRepo.GetById(id)
	if err != nil {
		return domain.User{}, nil, user.NotFoundError
	}

	return createdUser, sess, nil
}

func isLoginFormValid(uForm *domain.User) bool {
	return uForm.Email != "" && uForm.InputPassword != ""
}

func (uc *userUsecase) Login(u *domain.User) (domain.User, *session.Session, error) {
	if !isLoginFormValid(u) {
		return domain.User{}, nil, user.InvalidForm
	}
	foundUser, err := uc.userRepo.GetByEmail(u.Email)
	if err != nil {
		return domain.User{}, nil, user.NotFoundError
	}
	preparePassword(u)
	if foundUser.Password != u.Password {
		return domain.User{}, nil, user.InvalidCredentials
	}

	sess := &session.Session{UserID: foundUser.ID}
	err = session.Manager.Create(sess)
	if err != nil {
		return domain.User{}, nil, err
	}
	return foundUser, sess, nil
}

func (uc *userUsecase) Logout(sess *session.Session) (*session.Session, error) {
	err := session.Manager.Delete(sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
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
	return uc.userRepo.Delete(id)
}

func (uc *userUsecase) GetFavourites(id uint, sess *session.Session) ([]domain.Movie, error) {
	err := session.Manager.Check(sess)
	if err != nil || sess.UserID != id {
		return nil, err
	}

	return uc.userRepo.GetFavouritesByID(id)
}
