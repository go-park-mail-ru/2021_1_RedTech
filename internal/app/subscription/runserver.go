package subscription

import (
	"Redioteka/internal/pkg/database"
	subGRPC "Redioteka/internal/pkg/subscription/delivery/grpc"
	"Redioteka/internal/pkg/subscription/delivery/grpc/proto"
	"Redioteka/internal/pkg/subscription/repository"
	"Redioteka/internal/pkg/subscription/usecase"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func RunServer(addr string) {
	db := database.Connect()
	subRepo := repository.NewSubscriptionRepository(db)

	subUsecase := usecase.NewSubscriptionUsecase(subRepo)

	subGRPCHandler := subGRPC.NewSubscriptionHandler(subUsecase)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Log.Error(err)
	}
	server := grpc.NewServer()

	proto.RegisterSubscriptionServer(server, subGRPCHandler)

	log.Log.Info("starting server at " + addr)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		closeConnections(db)
		os.Exit(0)
	}()

	err = server.Serve(listen)
	if err != nil {
		log.Log.Error(err)
	}
}

func closeConnections(db *database.DBManager) {
	database.Disconnect(db)
	session.Destruct()
}
