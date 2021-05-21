package info

import (
	server2 "Redioteka/internal/app/server"
	_actorHandler "Redioteka/internal/pkg/actor/delivery/http"
	_actorRepository "Redioteka/internal/pkg/actor/repository"
	_actorUsecase "Redioteka/internal/pkg/actor/usecase"
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/middlewares"
	_movieHandler "Redioteka/internal/pkg/movie/delivery/http"
	_movieRepository "Redioteka/internal/pkg/movie/repository"
	_movieUsecase "Redioteka/internal/pkg/movie/usecase"
	_searchHandler "Redioteka/internal/pkg/search/delivery/http"
	_searchUsecase "Redioteka/internal/pkg/search/usecase"
	_userHandler "Redioteka/internal/pkg/user/delivery/http"
	"Redioteka/internal/pkg/user/repository"
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
	s.Use(middL.CSRFMiddleware)
	s.Use(middL.LoggingMiddleware)

	db := database.Connect()
	userRepo := repository.NewUserRepository(db)
	movieRepo := _movieRepository.NewMovieRepository(db)
	actorRepo := _actorRepository.NewActorRepository(db)
	avatarRepo := repository.NewS3AvatarRepository()

	userUsecase := _userUsecase.NewUserUsecase(userRepo, avatarRepo)
	movieUsecase := _movieUsecase.NewMovieUsecase(movieRepo)
	actorUsecase := _actorUsecase.NewActorUsecase(actorRepo)
	searchUsecase := _searchUsecase.NewSearchUsecase(movieRepo, actorRepo)

	_userHandler.NewUserHandlers(s, userUsecase)
	_movieHandler.NewMovieHandlers(s, movieUsecase)
	_actorHandler.NewActorHanlders(s, actorUsecase)
	_searchHandler.NewSearchHandlers(s, searchUsecase)

	// Static files
	fileRouter := r.PathPrefix("/static").Subrouter()
	server2.NewFileHandler(fileRouter)

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
