package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"

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

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := storage.RegisterElwinStorageHandlerFromEndpoint(ctx, gwmux, *storageEndpoint, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, strings.NewReader(swagger))
	})
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
