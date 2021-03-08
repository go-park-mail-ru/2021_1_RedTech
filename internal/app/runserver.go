package server

import (
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/user"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, "+
			"Content-Language, Content-Type")
		next.ServeHTTP(w, r)
	})
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, `{"msg": "Hello from server Redioteka"})`)
}

func RunServer(addr string) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	userApi := &user.Handler{}
	movieApi := &movie.Handler{}

	// Middleware
	s.Use(loggingMiddleware)
	s.Use(CORSMiddleware)

	// ===== Handlers start =====
	s.HandleFunc("/", handleRoot)

	// Users

	s.HandleFunc("/users/signup", userApi.Signup).Methods("POST", "OPTIONS")

	s.HandleFunc("/users/login", userApi.Login).Methods("POST", "OPTIONS")

	s.HandleFunc("/users/logout", userApi.Logout).Methods("GET", "OPTIONS")

	s.HandleFunc("/users/{id:[0-9]+}", userApi.Get).Methods("GET", "OPTIONS")

	s.HandleFunc("/me", userApi.Me).Methods("GET", "OPTIONS")

	s.HandleFunc("/users/{id}", userApi.Update).Methods("PATCH", "OPTIONS")

	s.HandleFunc("/users/{id}/avatar", userApi.Avatar)

	// Media

	s.HandleFunc("/media/movie/{id:[0-9]+}", movieApi.Get).Methods("GET", "OPTIONS")

	server := http.Server{
		Addr:    addr,
		Handler: s,
	}

	fmt.Println("starting server at", addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
