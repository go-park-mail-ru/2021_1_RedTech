package server

import (
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/user"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer(addr string) {
	r := mux.NewRouter()

	userApi := &user.UserHandler{}
	movieApi := &movie.MovieHandler{}

	// Users

	r.HandleFunc("/users/signup", userApi.Signup)

	r.HandleFunc("/users/login", userApi.Login)

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
