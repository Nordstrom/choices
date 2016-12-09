// Copyright 2016 Andrew O'Neill

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/foolusion/choices/storage/mongo"
	"github.com/foolusion/elwinprotos/storage"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	log.Println("Starting exp-store...")
	server, err := mongo.NewServer("localhost", "elwin")
	server.LoadExampleData()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	storage.RegisterElwinStorageServer(s, server)
	grpc_prometheus.Register(s)
	go func() {
		http.Handle("/metrics", prometheus.Handler())
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	log.Fatal(s.Serve(lis))
}
