package usecase

import (
	"Redioteka/internal/pkg/domain"
)

type searchUsecase struct {
	movieRepo domain.MovieRepository
	actorRepo domain.ActorRepository
}

func NewSearchUsecase(m domain.MovieRepository, a domain.ActorRepository) domain.SearchUsecase {
	return &searchUsecase{
		movieRepo: m,
		actorRepo: a,
	}
}

func (s *searchUsecase) Get(query string) (domain.Search, error) {
	movies, err := s.movieRepo.Search(query)
	if err != nil {
		return domain.Search{}, err
	}
	actors, err := s.actorRepo.Search(query)
	if err != nil {
		return domain.Search{}, err
	}
	res := domain.Search{
		Movies: movies,
		Actors: actors,
	}
	return res, nil
}
