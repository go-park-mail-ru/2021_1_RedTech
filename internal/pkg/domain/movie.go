package domain

type MovieType string

const (
	SeriesT MovieType = "series"
	MovieT  MovieType = "movie"
)

type Movie struct {
	ID          uint      `json:"id,omitempty" fake:"{number:1,100000}"`
	Rating      float32   `json:"rating,omitempty" fake:"{number:1,10}"`
	Title       string    `json:"title,omitempty" fake:"{sentence:3}"`
	Description string    `json:"description,omitempty" fake:"{sentence:25}"`
	Countries   []string  `json:"countries,omitempty"`
	IsFree      bool      `json:"is_free,omitempty"`
	Genres      []string  `json:"genres,omitempty"`
	Actors      []string  `json:"actors,omitempty"`
	Avatar      string    `json:"movie_avatar,omitempty"`
	Type        MovieType `json:"type,omitempty"`
	Year        string    `json:"year,omitempty"`
	Director    []string  `json:"director,omitempty"`
	Video       string    `json:"video_path,omitempty"`
}

const (
	FilterBoth = iota
	FilterFree
	FilterSubscription
)

func (m Movie) Stream() Movie {
	return Movie{Video: m.Video}
}

func (m Movie) Preview() Movie {
	return Movie{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Avatar:      m.Avatar,
	}
}

func (m Movie) Info() Movie {
	newM := m
	newM.Video = ""
	return newM
}

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
