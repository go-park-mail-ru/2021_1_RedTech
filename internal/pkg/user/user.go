package user

import (
	"net/http"
	"sync"
)

const hashLen = 32

type User struct {
	ID       uint          `json:"id"`
	Email    string        `json:"email"`
	Username string        `json:"username"`
	Password [hashLen]byte `json:"-"`
}

type Handler struct {
}

type usersData struct {
	sync.Mutex
	users map[string]*User
}

var data = usersData{
	users: make(map[string]*User),
}

func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Me(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Avatar(w http.ResponseWriter, r *http.Request) {
}
