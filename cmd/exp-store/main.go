package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/foolusion/choices/elwinstorage"
	"github.com/foolusion/choices/storage/mongo"
)

func main() {
	log.Println("Starting exp-store...")
	server, err := mongo.NewServer("localhost", "elwin")
	server.LoadExampleData()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	storage.RegisterElwinStorageServer(s, server)
	log.Fatal(s.Serve(lis))
}
