package domain

type AuthorizationUsecase interface {
	GetById(id uint) (User, error)
	Signup(u *User) (User, error)
	Login(u *User) (User, error)
	Update(u *User) error
	Delete(id uint) error
}

// authorization repository is equal to user repository
