package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	gw "github.com/foolusion/choices/elwinstorage"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	storageEndpoint = flag.String("storage_endpoint", "elwin-storage:80", "endpoint of elwin-storage")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterElwinStorageHandlerFromEndpoint(ctx, mux, *storageEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()
	log.Println("Starting json-gateway...")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
