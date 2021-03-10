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

var whiteListOrigin = map[string]struct{}{
	"http://localhost":           {},
	"http://redioteka.com":       {},
	"http://redioteka.com:3000":  {},
	"http://89.208.198.192:3000": {},
}

func CORSMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, found := whiteListOrigin[origin]; found {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			log.Printf("Request from unknown host: %s", origin)
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, "+
			"Content-Language, Content-Type, Content-Encoding")
		if r.Method == "OPTIONS" {
			return
		}
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
	r.HandleFunc("/", handleRoot)

	// Users

	s.HandleFunc("/users/signup", userApi.Signup).Methods("POST", "OPTIONS")

	s.HandleFunc("/users/login", userApi.Login).Methods("POST", "OPTIONS")

	s.HandleFunc("/users/logout", userApi.Logout).Methods("GET", "OPTIONS")

	s.HandleFunc("/users/{id:[0-9]+}", userApi.Get).Methods("GET", "OPTIONS")

	s.HandleFunc("/me", userApi.Me).Methods("GET", "OPTIONS")

	s.HandleFunc("/users/{id:[0-9]+}", userApi.Update).Methods("PATCH", "OPTIONS")

	s.HandleFunc("/users/{id:[0-9]+}/avatar", userApi.Avatar).Methods("POST", "PUT", "OPTIONS")

	// Media

	s.HandleFunc("/media/movie/{id:[0-9]+}", movieApi.Get).Methods("GET", "OPTIONS")

	// Static Files

	static := http.FileServer(http.Dir("./img"))
	r.PathPrefix("/static/movies/").Handler(http.StripPrefix("/static/", static))
	r.PathPrefix("/static/actors/").Handler(http.StripPrefix("/static/", static))
	r.PathPrefix("/static/users/").Handler(http.StripPrefix("/static/", static))

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
