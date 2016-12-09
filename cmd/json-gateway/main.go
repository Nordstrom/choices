package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	storageEndpoint = flag.String("storage_endpoint", "elwin-storage:80", "endpoint of elwin-storage")
	listenAddress   = flag.String("listen_address", ":8080", "address to listen on")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := storage.RegisterElwinStorageHandlerFromEndpoint(ctx, mux, *storageEndpoint, opts)
	if err != nil {
		return err
	}

	log.Printf("Listening on %s", *listenAddress)
	return http.ListenAndServe(*listenAddress, mux)
}

func main() {
	flag.Parse()
	log.Println("Starting json-gateway...")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
