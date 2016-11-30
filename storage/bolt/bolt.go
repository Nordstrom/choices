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

package bolt

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"golang.org/x/net/context"

	"github.com/boltdb/bolt"
	storage "github.com/foolusion/choices/elwinstorage"
	"github.com/golang/protobuf/proto"
)

var (
	environmentStaging    = []byte("staging")
	environmentProduction = []byte("production")
)

type Server struct {
	db *bolt.DB
}

func NewServer(file string) (*Server, error) {
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
	return &Server{db: db}, nil
}

func (s *Server) Close() error {
	return s.db.Close()
}

// All returns all the namespaces for a given environment.
func (s *Server) All(ctx context.Context, r *storage.AllRequest) (*storage.AllReply, error) {
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
func (s *Server) Create(ctx context.Context, r *storage.CreateRequest) (*storage.CreateReply, error) {
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
	cr := &storage.CreateReply{}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(env).Put([]byte(ns.Name), pns)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	cr.Namespace = ns
	return cr, nil
}

// Read returns the namespace matching the supplied name from the given
// environment.
func (s *Server) Read(ctx context.Context, r *storage.ReadRequest) (*storage.ReadReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	if len(r.Name) == 0 {
		return nil, fmt.Errorf("name is empty")
	}

	ns := storage.Namespace{}
	resp := &storage.ReadReply{}
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
	resp.Namespace = &ns
	return resp, nil
}

// Update replaces the namespace in the given environment with the namespace
// supplied.
func (s *Server) Update(ctx context.Context, r *storage.UpdateRequest) (*storage.UpdateReply, error) {
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
	resp := &storage.UpdateReply{}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(env).Put([]byte(ns.Name), pns)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	resp.Namespace = ns
	return resp, nil
}

// Delete deletes the namespace from the given environment.
func (s *Server) Delete(ctx context.Context, r *storage.DeleteRequest) (*storage.DeleteReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	if len(r.Name) == 0 {
		return nil, fmt.Errorf("name is empty")
	}

	ns := storage.Namespace{}
	resp := &storage.DeleteReply{}
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

	resp.Namespace = &ns
	return resp, nil
}

func envFromStorageRequest(e storage.Environment) []byte {
	switch e {
	case storage.Environment_Staging:
		return environmentStaging
	case storage.Environment_Production:
		return environmentProduction
	default:
		return environmentStaging
	}
	return environmentStaging
}
