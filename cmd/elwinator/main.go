package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/foolusion/choices"
	storage "github.com/foolusion/choices/elwinstorage"
	"google.golang.org/grpc"
)

const (
	storageAddr = "localhost:8080"
)

var (
	conn *grpc.ClientConn
	esc  storage.ElwinStorageClient
)

func init() {
	http.HandleFunc("/api/v1/experiments", experimentsHandler)
}

func main() {
	log.Println("Starting elwinator...")
	var err error
	conn, err = grpc.Dial(storageAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not dial grpc storage: %s", err)
	}
	esc = storage.NewElwinStorageClient(conn)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func experimentsHandler(w http.ResponseWriter, r *http.Request) {
	ar := &storage.AllRequest{
		Environment: storage.Environment_Staging,
	}
	allResp, err := esc.All(context.TODO(), ar)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cns []choices.Namespace
	for _, ns := range allResp.GetNamespaces() {
		n, err := choices.FromNamespace(ns)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cns = append(cns, n)
	}
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = enc.Encode(cns)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "could not encode json: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
