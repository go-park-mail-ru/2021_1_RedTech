package auth

import (
	"Redioteka/internal/constants"
	handler "Redioteka/internal/pkg/authorization/delivery/grpc"
	pb "Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/user/repository"
	_userUsecase "Redioteka/internal/pkg/user/usecase"
	"Redioteka/internal/pkg/utils/session"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func RunServer(addr string) {
	// All about data
	db := database.Connect(constants.DBUser, constants.DBPassword,
		constants.DBHost, constants.DBPort, constants.DBName)
	userRepo := repository.NewUserRepository(db)
	avatarRepo := repository.NewS3AvatarRepository()

	userUsecase := _userUsecase.NewUserUsecase(userRepo, avatarRepo)
	authHandler := handler.NewAuthorizationHandler(userUsecase, session.Manager)

	// All about grpc server
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	server := grpc.NewServer()

	pb.RegisterAuthorizationServer(server, authHandler)

	log.Print("starting server at", addr)

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
		log.Fatalln("Serve auth error: ", err)
	}
}

func closeConnections(db *database.DBManager) {
	database.Disconnect(db)
	session.Destruct()
}
