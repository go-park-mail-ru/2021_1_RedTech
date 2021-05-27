package domain

type Actor struct {
	ID        uint    `json:"id,omitempty"`
	FirstName string  `json:"first_name,omitempty"`
	LastName  string  `json:"last_name,omitempty"`
	Born      string  `json:"born,omitempty"`
	Avatar    string  `json:"actor_avatar,omitempty"`
	Movies    []Movie `json:"movies,omitempty"`
}

//go:generate mockgen -destination=../actor/repository/mock/mock_repo.go -package=mock Redioteka/internal/pkg/domain ActorRepository
type ActorRepository interface {
	GetById(id uint) (Actor, error)
	GetByMovie(movieID uint) ([]*Actor, error)
	Search(query string) ([]Actor, error)
}

//go:generate mockgen -destination=../actor/usecase/mock/mock_usecase.go -package=mock Redioteka/internal/pkg/domain ActorUsecase
type ActorUsecase interface {
	GetById(id uint) (Actor, error)
	Search(query string) ([]Actor, error)
}
