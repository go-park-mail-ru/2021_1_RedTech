package user

import "net/http"

type User struct {
ID       uint   `json:"id"`
Email    string `json:"email"`
Username string `json:"username"`
Password string `json:"password"`
}

type UserHandler struct {
}

func (api *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
}

func (api *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
}

func (api *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
}

func (api *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
}

func (api *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
}

func (api *UserHandler) Avatar(w http.ResponseWriter, r *http.Request) {
}
