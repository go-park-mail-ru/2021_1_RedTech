package auth

import (
	"Redioteka/internal/pkg/database"
	grpc3 "Redioteka/internal/pkg/user/delivery/grpc"
	pb "Redioteka/internal/pkg/user/delivery/grpc/proto"
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
	db := database.Connect()
	userRepo := repository.NewUserRepository(db)
	avatarRepo := repository.NewS3AvatarRepository()

	userUsecase := _userUsecase.NewUserUsecase(userRepo, avatarRepo)
	authHandler := grpc3.NewAuthorizationHandler(userUsecase)

	// All about grpc server
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	server := grpc.NewServer()

	pb.RegisterAuthorizationServer(server, authHandler)

	log.Print("starting server at :8081")

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
