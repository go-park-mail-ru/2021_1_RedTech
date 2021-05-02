package domain

type Actor struct {
	ID      uint     `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Surname string   `json:"surname,omitempty"`
	Born    string   `json:"born,omitempty"`
	Avatar  string   `json:"actor_avatar,omitempty"`
	Movies  []string `json:"movies,omitempty"`
}

type ActorRepository interface {
	GetById(id uint) (Actor, error)
}

type ActorUsecase interface {
	GetById(id uint) (Actor, error)
}
