package main

import (
	"Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	//stream.RunServer(":8083")
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewAuthorizationClient(conn)
}
