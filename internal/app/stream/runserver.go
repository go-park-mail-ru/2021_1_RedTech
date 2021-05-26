package stream

import (
	_authorizationProto "Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/config"
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/middlewares"
	_streamHandlers "Redioteka/internal/pkg/stream/delivery/http"
	_streamRepository "Redioteka/internal/pkg/stream/repository"
	_streamUsecase "Redioteka/internal/pkg/stream/usecase"
	"Redioteka/internal/pkg/utils/fileserver"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func RunServer(addr string) {
	authConn, err := grpc.Dial(config.Get().Auth.Host+config.Get().Auth.Port,
		grpc.WithInsecure())
	if err != nil {
		log.Log.Warn(fmt.Sprint("Can't connect to grpc ", err))
		return
	}
	defer authConn.Close()
	authClient := _authorizationProto.NewAuthorizationClient(authConn)
	sessionManager := session.NewGrpcSession(authClient)
	log.Log.Info(fmt.Sprint("Successfully connected to authorization server ",
		config.Get().Auth.Host+config.Get().Auth.Port))

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	middL := middlewares.InitMiddleware()
	r.Use(middL.PanicRecoverMiddleware)
	s.Use(middL.MetricsMiddleware)
	s.Use(middL.CORSMiddleware)
	s.Use(middL.CSRFMiddleware)
	s.Use(middL.LoggingMiddleware)

	db := database.Connect()
	streamRepo := _streamRepository.NewStreamRepository(db)

	streamUsecase := _streamUsecase.NewStreamUsecase(streamRepo)
	_streamHandlers.NewStreamHandlers(s, streamUsecase, sessionManager)

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

	err = server.ListenAndServe()
	if err != nil {
		log.Log.Error(err)
	}
}

func closeConnections(db *database.DBManager) {
	database.Disconnect(db)
	session.Destruct()
}
