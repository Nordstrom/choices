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
