package domain

type Actor struct {
	ID       uint    `json:"id,omitempty"`
	FirstName string  `json:"first_name,omitempty"`
	LastName string  `json:"last_name,omitempty"`
	Born     string  `json:"born,omitempty"`
	Avatar   string  `json:"actor_avatar,omitempty"`
	Movies   []Movie `json:"movies,omitempty"`
}

type ActorRepository interface {
	GetById(id uint) (Actor, error)
}

type ActorUsecase interface {
	GetById(id uint) (Actor, error)
}
