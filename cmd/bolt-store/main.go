// Copyright 2016 Andrew O'Neill, Nordstrom

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
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/context"

	"github.com/boltdb/bolt"
	"github.com/foolusion/elwinprotos/storage"
	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func bind(s []string) error {
	if len(s) == 0 {
		return nil
	}
	if err := viper.BindEnv(s[0]); err != nil {
		return err
	}
	return bind(s[1:])
}

func main() {
	log.Println("Starting bolt-store...")

	viper.SetDefault("db_file", "test.db")
	viper.SetDefault("listen_address", ":8080")
	viper.SetDefault("metrics_address", ":8081")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/elwin/bolt-store")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config: %v", err)
	}

	viper.SetEnvPrefix("bolt_store")
	if err := bind([]string{
		"db_file",
		"listen_address",
		"metrics_address",
	}); err != nil {
		log.Fatal(err)
	}

	server, err := newServer(viper.GetString("db_file"))

	log.Printf("lisening for grpc on %q", viper.GetString("listen_address"))
	lis, err := net.Listen("tcp", viper.GetString("listen_address"))
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
		log.Printf("listening for /metrics on %q", viper.GetString("metrics_address"))
		log.Fatal(http.ListenAndServe(viper.GetString("metrics_address"), nil))
	}()

	log.Fatal(s.Serve(lis))
}

var (
	environmentStaging    = []byte("staging")
	environmentProduction = []byte("production")
)

type server struct {
	db *bolt.DB
}

func newServer(file string) (*server, error) {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket(environmentStaging); err != nil {
			if err != bolt.ErrBucketExists {
				return err
			}
		}
		if _, err := tx.CreateBucket(environmentProduction); err != nil {
			if err != bolt.ErrBucketExists {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &server{db: db}, nil
}

func (s *server) Close() error {
	return s.db.Close()
}

// All returns all the namespaces for a given environment.
func (s *server) All(ctx context.Context, r *storage.AllRequest) (*storage.AllReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	ar := &storage.AllReply{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(env).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ns storage.Namespace
			if err := proto.Unmarshal(v, &ns); err != nil {
				return err
			}
			ar.Namespaces = append(ar.Namespaces, &ns)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ar, nil
}

// Create creates a namespace in the given environment.
func (s *server) Create(ctx context.Context, r *storage.CreateRequest) (*storage.CreateReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	ns := r.Namespace
	if ns == nil {
		return nil, fmt.Errorf("namespace is nil")
	}

	pns, err := proto.Marshal(ns)
	if err != nil {
		return nil, err
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(env).Put([]byte(ns.Name), pns)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.CreateReply{Namespace: ns}, nil
}

// Read returns the namespace matching the supplied name from the given
// environment.
func (s *server) Read(ctx context.Context, r *storage.ReadRequest) (*storage.ReadReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	if len(r.Name) == 0 {
		return nil, fmt.Errorf("name is empty")
	}

	ns := storage.Namespace{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket(env).Get([]byte(r.Name))
		if buf == nil {
			return grpc.Errorf(codes.NotFound, "key not found")
		}
		if err := proto.Unmarshal(buf, &ns); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.ReadReply{Namespace: &ns}, nil
}

// Update replaces the namespace in the given environment with the namespace
// supplied.
func (s *server) Update(ctx context.Context, r *storage.UpdateRequest) (*storage.UpdateReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	ns := r.GetNamespace()
	if ns == nil {
		return nil, fmt.Errorf("namespace is nil")
	}

	pns, err := proto.Marshal(ns)
	if err != nil {
		return nil, err
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(env).Put([]byte(ns.Name), pns)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.UpdateReply{Namespace: ns}, nil
}

// Delete deletes the namespace from the given environment.
func (s *server) Delete(ctx context.Context, r *storage.DeleteRequest) (*storage.DeleteReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	if len(r.Name) == 0 {
		return nil, fmt.Errorf("name is empty")
	}

	ns := storage.Namespace{}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		buf := tx.Bucket(env).Get([]byte(r.Name))
		if buf == nil {
			return grpc.Errorf(codes.NotFound, "key not found")
		}
		if err := proto.Unmarshal(buf, &ns); err != nil {
			return err
		}
		return tx.Bucket(env).Delete([]byte(r.Name))
	}); err != nil {
		return nil, err
	}
	return &storage.DeleteReply{Namespace: &ns}, nil
}

func (s *server) ExperimentIntake(ctx context.Context, r *storage.ExperimentIntakeRequest) (*storage.ExperimentIntakeReply, error) {
	return &storage.ExperimentIntakeReply{}, nil
}

func envFromStorageRequest(e storage.Environment) []byte {
	switch e {
	case storage.Staging:
		return environmentStaging
	case storage.Production:
		return environmentProduction
	default:
		return environmentStaging
	}
}
