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

	userApi := &user.Handler{}
	movieApi := &movie.Handler{}

	// Middleware
	r.Use(loggingMiddleware)
	r.Use(CORSMiddleware)

	// ===== Handlers start =====
	r.HandleFunc("/", handleRoot)

	// Users

	r.HandleFunc("/users/signup", userApi.Signup).Methods("POST")

	r.HandleFunc("/users/login", userApi.Login).Methods("POST")

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
