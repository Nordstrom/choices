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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, strings.NewReader(swagger))
	})
	mux.Handle("/", allowCORS(gwmux))

	log.Printf("Listening on %s", *listenAddress)
	return http.ListenAndServe(*listenAddress, mux)
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	return
}

func main() {
	flag.Parse()
	log.Println("Starting json-gateway...")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
