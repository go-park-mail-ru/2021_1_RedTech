package server

import (
	"Redioteka/internal/pkg/middlewares"
	_movieHandler "Redioteka/internal/pkg/movie/delivery/http"
	_movieRepository "Redioteka/internal/pkg/movie/repository"
	_movieUsecase "Redioteka/internal/pkg/movie/usecase"
	_userHandler "Redioteka/internal/pkg/user/delivery/http"
	_userRepository "Redioteka/internal/pkg/user/repository"
	_userUsecase "Redioteka/internal/pkg/user/usecase"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RunServer(addr string) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	middL := middlewares.InitMiddleware()
	r.Use(middL.PanicRecoverMiddleware)
	s.Use(middL.CORSMiddleware)
	s.Use(middL.LoggingMiddleware)

	userRepo := _userRepository.NewMapUserRepository()
	movieRepo := _movieRepository.NewMapMovieRepository()

	userUsecase := _userUsecase.NewUserUsecase(userRepo)
	movieUsecase := _movieUsecase.NewMovieUsecase(movieRepo)

	_userHandler.NewUserHandlers(s, userUsecase)
	_movieHandler.NewMovieHandlers(s, movieUsecase)

	// Static files
	fileRouter := r.PathPrefix("/static").Subrouter()
	NewFileHandler(fileRouter)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Log.Debug(fmt.Sprint("starting server at ", addr))

	err := server.ListenAndServe()
	if err != nil {
		log.Log.Error(err)
	}
}
