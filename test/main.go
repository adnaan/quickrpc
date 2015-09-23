package main

import (
	"log"
	"net"

	"github.com/adnaan/quickrpc/test/todo"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	s := grpc.NewServer()
	todo.RegisterTodoServiceServer(s, &todo.Server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.Serve(lis)
}
