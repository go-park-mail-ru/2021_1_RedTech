package server

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/middlewares"
	_movieHandler "Redioteka/internal/pkg/movie/delivery/http"
	_movieRepository "Redioteka/internal/pkg/movie/repository"
	_movieUsecase "Redioteka/internal/pkg/movie/usecase"
	_userHandler "Redioteka/internal/pkg/user/delivery/http"
	_userRepository "Redioteka/internal/pkg/user/repository"
	_userUsecase "Redioteka/internal/pkg/user/usecase"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func RunServer(addr string) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	middL := middlewares.InitMiddleware()
	r.Use(middL.PanicRecoverMiddleware)
	s.Use(middL.CORSMiddleware)
	s.Use(middL.LoggingMiddleware)

	db := database.Connect()
	userRepo := _userRepository.NewUserRepository(db)
	movieRepo := _movieRepository.NewMovieRepository(db)

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

	log.Log.Info("starting server at " + addr)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		closeConnections(db)
		os.Exit(0)
	}()

	err := server.ListenAndServe()
	if err != nil {
		log.Log.Error(err)
	}
}

func closeConnections(db *database.DBManager) {
	database.Disconnect(db)
	session.Destruct()
}
