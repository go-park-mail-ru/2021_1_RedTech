package domain

const hashLen = 32

type User struct {
	ID       uint          `json:"id"`
	Email    string        `json:"email,omitempty"`
	Username string        `json:"username,omitempty"`
	Password [hashLen]byte `json:"-"`
	Avatar   string        `json:"avatar,omitempty"`
}

func (u User) private() User {
	return User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

func (u User) public() User {
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
	Store(user *User) error
	Delete(id uint) error
}

type UserUsecase interface {
	GetById(id uint) (User, error)
	GetByEmail(email string) (User, error)
	Update(user *User) error
	Delete(id uint) error
}

type MovieType string

const (
	series MovieType = "series"
	movie  MovieType = "movie"
)

type Movie struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rating      int       `json:"rating"`
	Countries   []string  `json:"countries"`
	IsFree      bool      `json:"is_free"`
	Genres      []string  `json:"genres"`
	Actors      []string  `json:"actors"`
	Avatar      string    `json:"movie_avatar,omitempty"`
	Type        MovieType `json:"type"`
	Year        string    `json:"year"`
	Director    []string  `json:"director"`
}

type MovieRepository interface {
	GetById(id uint) (Movie, error)
	Delete(id uint) error
}

type MovieUsecase interface {
	GetById(id uint) (Movie, error)
	Delete(id uint) error
}
