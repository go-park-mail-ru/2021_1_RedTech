package auth

import (
	handler "Redioteka/internal/pkg/authorization/delivery/grpc"
	pb "Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	_authUsecase "Redioteka/internal/pkg/authorization/usecase"
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/user/repository"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func RunServer(addr string) {
	// All about data
	db := database.Connect()
	userRepo := repository.NewUserRepository(db)

	userUsecase := _authUsecase.NewAuthorizationUsecase(userRepo)
	authHandler := handler.NewAuthorizationHandler(userUsecase, session.Manager)

	// All about grpc server
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Log.Error(err)
		return
	}
	server := grpc.NewServer()

	pb.RegisterAuthorizationServer(server, authHandler)

	log.Log.Info("starting auth grpc server at " + addr)

	// All about server start
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		closeConnections(db)
		os.Exit(0)
	}()

	err = server.Serve(lis)
	if err != nil {
		log.Log.Error(err)
		return
	}
}

func closeConnections(db *database.DBManager) {
	database.Disconnect(db)
	session.Destruct()
}
