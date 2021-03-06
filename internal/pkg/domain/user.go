package domain

import (
	"Redioteka/internal/pkg/utils/session"
	"io"
	"time"
)

type User struct {
	ID                   uint   `json:"id,omitempty"`
	Email                string `json:"email,omitempty"`
	Username             string `json:"username,omitempty"`
	Password             []byte `json:"-"`
	Avatar               string `json:"avatar,omitempty"`
	InputPassword        string `json:"password,omitempty"`
	ConfirmInputPassword string `json:"confirm_password,omitempty"`
	IsSubscriber         bool   `json:"is_sub"`
}

type UserFavourites struct {
	Favourites []Movie `json:"favourites"`
}

func (u User) Private() User {
	return User{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		Avatar:       u.Avatar,
		IsSubscriber: u.IsSubscriber,
	}
}

func (u User) Public() User {
	return User{
		ID:           u.ID,
		Username:     u.Username,
		Avatar:       u.Avatar,
		IsSubscriber: u.IsSubscriber,
	}
}

//go:generate mockgen -destination=../user/repository/mock/mock_repo.go -package=mock Redioteka/internal/pkg/domain UserRepository
type UserRepository interface {
	GetById(id uint) (User, error)
	GetByEmail(email string) (User, error)
	Update(user *User) error
	Store(user *User) (uint, error)
	Delete(id uint) error
	GetFavouritesByID(id uint) ([]Movie, error)
	CheckSub(id uint) time.Time
}

type AvatarRepository interface {
	UploadAvatar(reader io.Reader, path, ext string) (string, error)
}

//go:generate mockgen -destination=../user/usecase/mock/mock_usecase.go -package=mock Redioteka/internal/pkg/domain UserUsecase
type UserUsecase interface {
	GetById(id uint) (User, error)
	Signup(u *User) (User, error)
	Login(u *User) (User, error)
	Logout(sess *session.Session) error
	Update(u *User) error
	Delete(id uint) error
	GetFavourites(id uint, sess *session.Session) ([]Movie, error)
	UploadAvatar(reader io.Reader, path, ext string) (string, error)
}
