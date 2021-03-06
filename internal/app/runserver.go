package server

import (
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/user"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func RunServer(addr string) {
	r := mux.NewRouter()

	userApi := &user.Handler{}
	movieApi := &movie.Handler{}

	// Users

	r.HandleFunc("/users/signup", userApi.Signup)

	r.HandleFunc("/users/login", userApi.Login)

	r.HandleFunc("/users/logout", userApi.Logout)

	r.HandleFunc("/users/{id}", userApi.Login)

	r.HandleFunc("/users/{id}/avatar", userApi.Avatar)

	r.HandleFunc("/me", userApi.Me)

	// Media

	r.HandleFunc("/movie/{id}", movieApi.Get)

	// Middleware
	r.Use(loggingMiddleware)

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
