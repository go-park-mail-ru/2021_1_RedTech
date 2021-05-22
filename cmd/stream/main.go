package main

import (
	proto2 "Redioteka/internal/pkg/authorization/delivery/grpc/proto"
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
	c := proto2.NewAuthorizationClient(conn)

	// Contact the server and print out its response.
	ctx := context.TODO()
	r, err := c.CreateSession(ctx, &proto2.CreateSessionParams{
		UserId: 1,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetCookie())
}
