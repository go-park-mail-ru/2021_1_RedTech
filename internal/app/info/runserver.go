package info

import (
	"Redioteka/internal/constants"
	_actorHandler "Redioteka/internal/pkg/actor/delivery/http"
	_actorRepository "Redioteka/internal/pkg/actor/repository"
	_actorUsecase "Redioteka/internal/pkg/actor/usecase"
	_authorizationProto "Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/middlewares"
	_movieHandler "Redioteka/internal/pkg/movie/delivery/http"
	_movieRepository "Redioteka/internal/pkg/movie/repository"
	_movieUsecase "Redioteka/internal/pkg/movie/usecase"
	_searchHandler "Redioteka/internal/pkg/search/delivery/http"
	_searchUsecase "Redioteka/internal/pkg/search/usecase"
	_subscriptionProto "Redioteka/internal/pkg/subscription/delivery/grpc/proto"
	_subscriptionHandler "Redioteka/internal/pkg/subscription/delivery/http"
	_userHandler "Redioteka/internal/pkg/user/delivery/http"
	"Redioteka/internal/pkg/user/repository"
	_userUsecase "Redioteka/internal/pkg/user/usecase/grpc_usecase"
	"Redioteka/internal/pkg/utils/fileserver"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func RunServer(addr string) {
	// GRPC connecting
	conn, err := grpc.Dial(constants.AuthServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Log.Error(err)
	}
	defer conn.Close()
	authClient := _authorizationProto.NewAuthorizationClient(conn)

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	middlewares.RegisterMetrics()

	middL := middlewares.InitMiddleware()
	r.Use(middL.PanicRecoverMiddleware)
	s.Use(middL.MetricsMiddleware)
	s.Use(middL.CORSMiddleware)
	s.Use(middL.CSRFMiddleware)
	s.Use(middL.LoggingMiddleware)

	db := database.Connect(constants.DBUser, constants.DBPassword,
		constants.DBHost, constants.DBPort, constants.DBName)
	sessionManager := session.NewGrpcSession(authClient)
	userRepo := repository.NewUserRepository(db)
	movieRepo := _movieRepository.NewMovieRepository(db)
	actorRepo := _actorRepository.NewActorRepository(db)
	avatarRepo := repository.NewS3AvatarRepository()

	userUsecase := _userUsecase.NewGrpcUserUsecase(userRepo, avatarRepo, authClient, sessionManager)
	movieUsecase := _movieUsecase.NewMovieUsecase(movieRepo, userRepo, actorRepo)
	actorUsecase := _actorUsecase.NewActorUsecase(actorRepo)
	searchUsecase := _searchUsecase.NewSearchUsecase(movieRepo, actorRepo)

	_userHandler.NewUserHandlers(s, userUsecase)
	_movieHandler.NewMovieHandlers(s, movieUsecase)
	_actorHandler.NewActorHandlers(s, actorUsecase)
	_searchHandler.NewSearchHandlers(s, searchUsecase)

	grpcConn, err := grpc.Dial(
		"subscription:8084",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Log.Warn("cant connect to grpc")
		return
	}

	subClient := _subscriptionProto.NewSubscriptionClient(grpcConn)
	_subscriptionHandler.NewSubscriptionHandlers(r, subClient)
	r.Handle("/metrics", promhttp.Handler())

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
		closeConnections(db, grpcConn)
		os.Exit(0)
	}()

	err = server.ListenAndServe()
	if err != nil {
		log.Log.Error(err)
	}
}

func closeConnections(db *database.DBManager, grpcConn *grpc.ClientConn) {
	database.Disconnect(db)
	session.Destruct()
	grpcConn.Close()
}
