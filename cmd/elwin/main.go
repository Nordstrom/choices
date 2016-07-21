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
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/context"

	"github.com/foolusion/choices"
	"github.com/foolusion/choices/elwin"
	"google.golang.org/grpc"
)

func init() {
	http.HandleFunc("/", rootHandler)
}

func main() {
	log.Println("Starting elwin...")
	t1, err := choices.NewNamespace("t1", "test", []string{"userid"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	t1.Addexp(
		"uniform",
		[]choices.Param{{Name: "a", Value: &choices.Uniform{Choices: []string{"b", "c"}}}},
		128,
	)
	if err := choices.Addns(t1); err != nil {
		log.Fatalf("%v", err)
	}

	t2, err := choices.NewNamespace("t2", "test", []string{"userid"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	t2.Addexp(
		"weighted",
		[]choices.Param{{Name: "b", Value: &choices.Weighted{Choices: []string{"on", "off"}, Weights: []float64{2, 1}}}},
		128,
	)
	if err := choices.Addns(t2); err != nil {
		log.Fatalf("%v", err)
	}

	t3, err := choices.NewNamespace("t3", "test", []string{"userid"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	t3.Addexp(
		"halfSegments",
		[]choices.Param{{Name: "b", Value: &choices.Uniform{Choices: []string{"on"}}}},
		64,
	)
	if err := choices.Addns(t3); err != nil {
		log.Fatalf("%v", err)
	}

	t4, err := choices.NewNamespace("t4", "test", []string{"userid"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	t4.Addexp(
		"multi",
		[]choices.Param{
			{Name: "a", Value: &choices.Uniform{Choices: []string{"on", "off"}}},
			{Name: "b", Value: &choices.Weighted{Choices: []string{"up", "down"}, Weights: []float64{1, 2}}},
		},
		128,
	)
	if err := choices.Addns(t4); err != nil {
		log.Fatalf("%v", err)
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("main: failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	elwin.RegisterElwinServer(grpcServer, &elwinServer{})
	grpcServer.Serve(lis)
}

type elwinServer struct {
}

func (e *elwinServer) GetNamespaces(ctx context.Context, id *elwin.Identifier) (*elwin.Experiments, error) {
	if id == nil {
		return nil, fmt.Errorf("GetNamespaces: no Identifier recieved")
	}

	return choices.Namespaces(nil, id.TeamID, id.UserID)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	resp, err := choices.Namespaces(nil, "test", r.Form)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	enc := json.NewEncoder(w)
	enc.Encode(*resp)
}
