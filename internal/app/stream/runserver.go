package stream

import (
	"Redioteka/internal/constants"
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/middlewares"
	_streamHandlers "Redioteka/internal/pkg/stream/delivery/http"
	_streamRepository "Redioteka/internal/pkg/stream/repository"
	_streamUsecase "Redioteka/internal/pkg/stream/usecase"
	"Redioteka/internal/pkg/utils/fileserver"
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

	db := database.Connect(constants.DBUser, constants.DBPassword,
		constants.DBHost, constants.DBPort, constants.DBName)
	streamRepo := _streamRepository.NewStreamRepository(db)

	streamUsecase := _streamUsecase.NewStreamUsecase(streamRepo)
	_streamHandlers.NewStreamHandlers(s, streamUsecase, session.Manager)

	// Static files
	fileRouter := r.PathPrefix("/static").Subrouter()
	fileserver.NewFileHandler(fileRouter)

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
