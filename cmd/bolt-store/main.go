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
	"crypto/rand"
	"encoding/binary"
	"log"
	"net"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/foolusion/elwinprotos/storage"
	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/labels"
)

func bind(s ...string) error {
	if len(s) == 0 {
		return nil
	}
	if err := viper.BindEnv(s[0]); err != nil {
		return err
	}
	return bind(s[1:]...)
}

var (
	ErrNilRequest = errors.New("request is nil")
)

func main() {
	log.Println("Starting bolt-store...")

	viper.SetDefault("db_file", "test.db")
	viper.SetDefault("db_bucket", "dev")
	viper.SetDefault("listen_address", ":8080")
	viper.SetDefault("metrics_address", ":8081")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/elwin/bolt-store")
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("no config file found")
		default:
			log.Fatalf("could not read config: %v", err)
		}
	}

	viper.SetEnvPrefix("bolt_store")
	if err := bind(
		"db_file",
		"listen_address",
		"metrics_address",
	); err != nil {
		log.Fatal(err)
	}

	server, err := newServer(viper.GetString("db_file"), viper.GetString("db_bucket"))

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
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("listening for /metrics on %q", viper.GetString("metrics_address"))
		log.Fatal(http.ListenAndServe(viper.GetString("metrics_address"), nil))
	}()

	log.Fatal(s.Serve(lis))
}

type server struct {
	db     *bolt.DB
	bucket []byte
}

func newServer(file, bucket string) (*server, error) {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}

	if bucket == "" {
		return nil, errors.New("bucket is empty")
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket([]byte(bucket)); err != nil {
			if err != bolt.ErrBucketExists {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &server{db: db, bucket: []byte(bucket)}, nil
}

func (s *server) Close() error {
	return s.db.Close()
}

// List returns all the experiments that match a query.
func (s *server) List(ctx context.Context, r *storage.ListRequest) (*storage.ListReply, error) {
	if r == nil {
		return nil, ErrNilRequest
	}

	selector, err := labels.Parse(r.Query)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse query")
	}

	ar := &storage.ListReply{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(s.bucket).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var exp storage.Experiment
			if err := proto.Unmarshal(v, &exp); err != nil {
				return err
			}
			if selector.Matches(labels.Set(exp.Labels)) {
				ar.Experiments = append(ar.Experiments, &exp)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ar, nil
}

// Set creates an experiment in the given environment.
func (s *server) Set(ctx context.Context, r *storage.SetRequest) (*storage.SetReply, error) {
	if r == nil {
		return nil, ErrNilRequest
	}

	exp := r.Experiment
	if exp == nil {
		return nil, errors.New("experiment is nil")
	}

	if exp.Id == "" {
		// TODO: set exp.Id to a new generated id
		if name, err := randName(32); err != nil {
			errors.Wrap(err, "could not create random name")
		} else {
			exp.Id = name
		}
	}

	pexp, err := proto.Marshal(exp)
	if err != nil {
		return nil, err
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(s.bucket).Put([]byte(exp.Id), pexp)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.SetReply{Experiment: exp}, nil
}

func randName(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var str string
	b := make([]byte, 8*n)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "could not read from rand")
	}
	for i := 0; i < n; i++ {
		a := binary.BigEndian.Uint64(b[i*8:(i+1)*8]) % uint64(len(alphabet))
		str += alphabet[a : a+1]
	}
	return str, nil
}

// Get returns the experiment matching the supplied id from the given
// environment.
func (s *server) Get(ctx context.Context, r *storage.GetRequest) (*storage.GetReply, error) {
	if r == nil {
		return nil, ErrNilRequest
	}

	if r.Id == "" {
		return nil, errors.New("id is empty")
	}

	exp := storage.Experiment{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket(s.bucket).Get([]byte(r.Id))
		if buf == nil {
			return grpc.Errorf(codes.NotFound, "key not found")
		}
		if err := proto.Unmarshal(buf, &exp); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.GetReply{Experiment: &exp}, nil
}

// Remove deletes the experiment from the given environment.
func (s *server) Remove(ctx context.Context, r *storage.RemoveRequest) (*storage.RemoveReply, error) {
	if r == nil {
		return nil, ErrNilRequest
	}

	if r.Id == "" {
		return nil, errors.New("id is empty")
	}

	exp := storage.Experiment{}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		buf := tx.Bucket(s.bucket).Get([]byte(r.Id))
		if buf == nil {
			return grpc.Errorf(codes.NotFound, "key not found")
		}
		if err := proto.Unmarshal(buf, &exp); err != nil {
			return err
		}
		return tx.Bucket(s.bucket).Delete([]byte(r.Id))
	}); err != nil {
		return nil, err
	}
	return &storage.RemoveReply{Experiment: &exp}, nil
}
