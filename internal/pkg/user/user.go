package user

import "net/http"

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Handler struct {
}

func (api *Handler) Login(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Signup(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Logout(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Avatar(w http.ResponseWriter, r *http.Request) {
}


