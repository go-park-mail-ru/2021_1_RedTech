package server

import (
	"Redioteka/internal/app/middlewares"
	_movieHandler "Redioteka/internal/app/movie/delivery/http"
	_movieRepository "Redioteka/internal/app/movie/repository"
	_movieUsecase "Redioteka/internal/app/movie/usecase"
	_userHandler "Redioteka/internal/app/user/delivery/http"
	_userRepository "Redioteka/internal/app/user/repository"
	_userUsecase "Redioteka/internal/app/user/usecase"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer(addr string) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	middL := middlewares.InitMiddleware()
	s.Use(middL.CORSMiddleware)
	s.Use(middL.LoggingMiddleware)

	userRepo := _userRepository.NewMapUserRepository()
	movieRepo := _movieRepository.NewMapMovieRepository()

	userUsecase := _userUsecase.NewUserUsecase(userRepo)
	movieUsecase := _movieUsecase.NewMovieUsecase(movieRepo)

	_userHandler.NewUserHandlers(s, userUsecase)
	_movieHandler.NewMovieHandlers(s, movieUsecase)

	// Static files
	static := http.FileServer(http.Dir("./img"))
	r.PathPrefix("/static/movies/").Handler(http.StripPrefix("/static/", static))
	r.PathPrefix("/static/actors/").Handler(http.StripPrefix("/static/", static))
	r.PathPrefix("/static/users/").Handler(http.StripPrefix("/static/", static))

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Println("starting server at ", addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
