package domain

const hashLen = 32

type User struct {
	ID                   uint          `json:"id,omitempty"`
	Email                string        `json:"email,omitempty"`
	Username             string        `json:"username,omitempty"`
	Password             [hashLen]byte `json:"-"`
	Avatar               string        `json:"avatar,omitempty"`
	InputPassword        string        `json:"password,omitempty"`
	ConfirmInputPassword string        `json:"confirm_password,omitempty"`
}

func (u User) Private() User {
	return User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

func (u User) Public() User {
	return User{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

type UserRepository interface {
	GetById(id uint) (User, error)
	GetByEmail(email string) (User, error)
	Update(user *User) error
	Store(user *User) (uint, error)
	Delete(id uint) error
}

type UserUsecase interface {
	GetById(id uint) (User, error)
	GetCurrent() (User, error)
	Signup(u *User) (User, error)
	Login(u *User) (User, error)
	Logout(u *User) error
	Update(u *User) error
	Delete(id uint) error
}