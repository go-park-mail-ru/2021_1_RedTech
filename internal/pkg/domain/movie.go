package domain

import "Redioteka/internal/pkg/utils/session"

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
	IsFree      bool      `json:"is_free"`
	Genres      []string  `json:"genres,omitempty"`
	Actors      []string  `json:"actors,omitempty"`
	ActorIds    []uint    `json:"actor_ids,omitempty"`
	Avatar      string    `json:"movie_avatar,omitempty"`
	Type        MovieType `json:"type,omitempty"`
	Year        string    `json:"year,omitempty"`
	Director    []string  `json:"director,omitempty"`
	Favourite   int       `json:"is_fav,omitempty"`
	Vote        int       `json:"is_vote,omitempty"`
	Series      []uint    `json:"series_list,omitempty"`
}

type Genre struct {
	Name     string `json:"name"`
	LabelRus string `json:"label_rus"`
	Image    string `json:"image"`
}

type Stream struct {
	Video  string `json:"video_path,omitempty"`
	Season int    `json:"season,omitempty"`
	Series int    `json:"series,omitempty"`
}

const (
	Like    = 1
	Dislike = -1
)

func (m Movie) Preview() Movie {
	return Movie{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Avatar:      m.Avatar,
	}
}

const (
	FilterBoth = iota
	FilterFree
	FilterSubscription
)

const (
	NoneOrder = iota
	RatingOrder
	DateOrder
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
	Order     int       `schema:"-"`
}

//go:generate mockgen -destination=../movie/repository/mock/mock_repo.go -package=mock Redioteka/internal/pkg/domain MovieRepository
type MovieRepository interface {
	GetById(id uint) (Movie, error)
	AddFavouriteByID(movieID, userID uint) error
	RemoveFavouriteByID(movieID, userID uint) error
	CheckFavouriteByID(movieID, userID uint) error
	CheckVoteByID(movieID, userID uint) int
	GetByFilter(filter MovieFilter) ([]Movie, error)
	GetGenres() ([]Genre, error)
	GetSeriesList(id uint) ([]uint, error)
	GetStream(id uint) ([]Stream, error)
	Like(userId, movieId uint) error
	Dislike(userId, movieId uint) error
}

//go:generate mockgen -destination=../movie/usecase/mock/mock_usecase.go -package=mock Redioteka/internal/pkg/domain MovieUsecase
type MovieUsecase interface {
	GetByID(id uint, sess *session.Session) (Movie, error)
	AddFavourite(id uint, sess *session.Session) error
	RemoveFavourite(id uint, sess *session.Session) error
	GetByFilter(filter MovieFilter) ([]Movie, error)
	GetGenres() ([]Genre, error)
	GetStream(id uint) ([]Stream, error)
	Like(userId, movieId uint) error
	Dislike(userId, movieId uint) error
}
