package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/foolusion/elwinprotos/elwin"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	storageEndpoint = flag.String("storage_endpoint", "localhost:8080", "endpoint of elwin")
	listenAddress   = flag.String("listen_address", ":8081", "address to listen on")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor)}
	cc, err := grpc.Dial(*storageEndpoint, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	e := &experiment{
		ElwinClient: elwin.NewElwinClient(cc),
	}

	grpc_prometheus.EnableClientHandlingTimeHistogram(grpc_prometheus.WithHistogramBuckets(prometheus.ExponentialBuckets(.0000000001, 10, 10)))
	mux := http.NewServeMux()
	mux.HandleFunc("/elwin/v1/experiments", e.experimentHandler)
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s", *listenAddress)
	return http.ListenAndServe(*listenAddress, mux)
}

type experiment struct {
	elwin.ElwinClient
}

func (e *experiment) experimentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	var data elwin.GetRequest
	if err := dec.Decode(&data); err != nil {
		http.Error(w, "could not parse json: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	resp, err := e.Get(ctx, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		http.Error(w, "could not marshal json", http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.Parse()
	log.Println("Starting json-gateway...")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
