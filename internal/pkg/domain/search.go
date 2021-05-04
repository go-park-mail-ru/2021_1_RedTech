package domain

type Search struct {
	Movies []Movie
	Actors []Actor
}

type SearchRepository interface {
	Get(query string) (Search, error)
}

type SearchUsecase interface {
	Get(query string) (Search, error)
}

