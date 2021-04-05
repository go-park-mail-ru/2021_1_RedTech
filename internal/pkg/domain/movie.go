package domain

type MovieType string

const (
	SeriesT MovieType = "series"
	MovieT  MovieType = "movie"
)

type Movie struct {
	ID          uint      `json:"id" fake:"{number:1,100000}"`
	Title       string    `json:"title" fake:"{sentence:3}"`
	Description string    `json:"description" fake:"{sentence:25}"`
	Rating      float32   `json:"rating" fake:"{number:1,10}"`
	Countries   []string  `json:"countries"`
	IsFree      bool      `json:"is_free"`
	Genres      []string  `json:"genres"`
	Actors      []string  `json:"actors"`
	Avatar      string    `json:"movie_avatar,omitempty"`
	Type        MovieType `json:"type"`
	Year        string    `json:"year"`
	Director    []string  `json:"director"`
	Video       string    `json:"video_path"`
}

const (
	FilterBoth = iota
	FilterFree
	FilterSubscription
)

type MovieFilter struct {
	MinRating float32   `schema:"min_rating"`
	Countries []string  `schema:"countries"`
	IsFree    int       `schema:"is_free"` // FilterFree | FilterSubscription | FilterBoth
	Genres    []string  `schema:"genres"`
	Actors    []string  `schema:"actors"`
	Type      MovieType `schema:"type"`
	Director  []string  `schema:"director"`
	Offset    int       `schema:"offset"`
	Limit     int       `schema:"limit"`
}

//go:generate mockgen -destination=../movie/repository/mock/mock_repo.go -package=mock Redioteka/internal/pkg/domain MovieRepository
type MovieRepository interface {
	GetById(id uint) (Movie, error)
	GetByFilter(filter MovieFilter) ([]Movie, error)
}

//go:generate mockgen -destination=../movie/usecase/mock/mock_usecase.go -package=mock Redioteka/internal/pkg/domain MovieUsecase
type MovieUsecase interface {
	GetById(id uint) (Movie, error)
	GetByFilter(filter MovieFilter) ([]Movie, error)
}
