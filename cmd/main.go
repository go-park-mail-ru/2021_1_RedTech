package main

/*
RT-35
1) Должен быть http-сервер

2) Выбрать роутер
*/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

type Movie struct {
}


type MovieHandler struct {
}

func (api *MovieHandler) Get(w http.ResponseWriter, r *http.Request) {
}

func runServer(addr string) {
	r := mux.NewRouter()

	userApi := &UserHandler{}
	movieApi := &MovieHandler{}

	// Users

	r.HandleFunc("/users/signup", userApi.Signup)

	r.HandleFunc("/users/login/", userApi.Login)

	r.HandleFunc("/users/logout", userApi.Logout)

	r.HandleFunc("/users/{id}", userApi.Login)

	r.HandleFunc("/users/{id}/avatar", userApi.Avatar)

	r.HandleFunc("/me", userApi.Me)

	// Media
	r.HandleFunc("/movie/{id}", movieApi.Get)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Println("starting server at", addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go runServer(":8081")
	runServer(":8080")
}
