package domain

type MovieType string

const (
	SeriesT MovieType = "series"
	MovieT  MovieType = "movie"
)

type Movie struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rating      float32   `json:"rating"`
	Countries   []string  `json:"countries,omitempty"`
	IsFree      bool      `json:"is_free"`
	Genres      []string  `json:"genres,omitempty"`
	Actors      []string  `json:"actors,omitempty"`
	Avatar      string    `json:"movie_avatar,omitempty"`
	Type        MovieType `json:"type,omitempty"`
	Year        string    `json:"year,omitempty"`
	Director    []string  `json:"director,omitempty"`
}

//go:generate mockgen -destination=../movie/repository/mock/mock_repo.go -package=mock Redioteka/internal/pkg/domain MovieRepository
type MovieRepository interface {
	GetById(id uint) (Movie, error)
}

//go:generate mockgen -destination=../movie/usecase/mock/mock_usecase.go -package=mock Redioteka/internal/pkg/domain MovieUsecase
type MovieUsecase interface {
	GetById(id uint) (Movie, error)
}
