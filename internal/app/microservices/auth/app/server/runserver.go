package server

import (
	grpc2 "Redioteka/internal/app/microservices/auth/delivery/grpc"
	"Redioteka/internal/app/microservices/auth/delivery/grpc/proto"
	_actorUsecase "Redioteka/internal/pkg/actor/usecase"
	"Redioteka/internal/pkg/database"
	_movieUsecase "Redioteka/internal/pkg/movie/usecase"
	_searchUsecase "Redioteka/internal/pkg/search/usecase"
	_avatarRepository "Redioteka/internal/pkg/user/repository"
	_userRepository "Redioteka/internal/pkg/user/repository"
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
	userRepo := _userRepository.NewUserRepository(db)
	avatarRepo := _avatarRepository.NewS3AvatarRepository()

	userUsecase := _userUsecase.NewUserUsecase(userRepo, avatarRepo)
	authHandler := grpc2.NewAuthorizationHandler(userRepo, )

	// All about grpc server
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listen port", err
	}
	server := grpc.NewServer()

	proto.RegisterAuthorizationServer(server)

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
