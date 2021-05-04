package domain

type Search struct {
	Movies []Movie `json:"movies,omitempty"`
	Actors []Actor `json:"actors,omitempty"`
}

type SearchUsecase interface {
	Get(query string) (Search, error)
}
