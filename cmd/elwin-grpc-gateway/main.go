package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/foolusion/elwinprotos/elwin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	storageEndpoint = flag.String("storage_endpoint", ":8083", "endpoint of elwin")
	listenAddress   = flag.String("listen_address", ":8085", "address to listen on")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := elwin.RegisterElwinHandlerFromEndpoint(ctx, gwmux, *storageEndpoint, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

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
