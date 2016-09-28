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

package mongo

import (
	"encoding/hex"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"golang.org/x/net/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	storage "github.com/foolusion/choices/elwinstorage"
	"github.com/foolusion/choices/storage/mongo/internal/types"
	"github.com/pkg/errors"
)

const (
	environmentStaging    = "staging"
	environmentProduction = "production"
)

type Server struct {
	sess *mgo.Session
	db   string
}

func NewServer(addr, db string) (*Server, error) {
	sess, err := mgo.Dial(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "could not dial %q", addr)
	}
	return &Server{sess: sess, db: db}, nil
}

func (s *Server) All(ctx context.Context, r *storage.AllRequest) (*storage.AllReply, error) {
	var env string
	switch {
	case r.Environment == storage.Environment_Staging:
		env = environmentStaging
	case r.Environment == storage.Environment_Production:
		env = environmentProduction
	default:
		return nil, grpc.Errorf(codes.InvalidArgument, "bad environment requested")
	}

	var results []types.Namespace
	err := s.sess.DB(s.db).C(env).Find(nil).All(&results)
	if err == mgo.ErrNotFound {
		return nil, grpc.Errorf(codes.NotFound, "could not find all: %s", err)
	} else if err != nil {
		return nil, grpc.Errorf(codes.Unknown, "could not decode data from mongo: %s", err)
	}

	resp, err := parseNamespaces(results)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "could not parse namespaces: %s", err)
	}

	return &storage.AllReply{
		Namespaces: resp,
	}, nil
}

func (s *Server) Create(ctx context.Context, r *storage.CreateRequest) (*storage.CreateReply, error) {
	log.Println("strarting create...")
	if r == nil || r.Namespace == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad request")
	}

	nsi := snToN(r.Namespace)
	var env string
	switch r.Environment {
	case storage.Environment_Staging:
		env = environmentStaging
	case storage.Environment_Production:
		env = environmentProduction
	default:
		return nil, grpc.Errorf(codes.InvalidArgument, "bad environment provided")
	}

	err := s.sess.DB(s.db).C(env).Insert(nsi)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to insert experiment into database: %s", err)
	}

	return &storage.CreateReply{Namespace: r.Namespace}, nil
}

func (s *Server) Read(ctx context.Context, r *storage.ReadRequest) (*storage.ReadReply, error) {
	log.Println("starting read...")
	var rr storage.ReadReply
	if r == nil || r.Name == "" {
		return &rr, grpc.Errorf(codes.InvalidArgument, "bad request")
	}

	var env string
	switch r.Environment {
	case storage.Environment_Staging:
		env = environmentStaging
	case storage.Environment_Production:
		env = environmentProduction
	default:
		return &rr, grpc.Errorf(codes.InvalidArgument, "bad environment provided")
	}

	ns, err := s.getNamespace(r.Name, env)
	switch {
	case err != nil:
		return &rr, err
	case ns == nil:
		return &rr, grpc.Errorf(codes.Internal, "nil response from getNamespace")
	}
	return &storage.ReadReply{Namespace: ns}, nil
}

func (s *Server) Update(ctx context.Context, r *storage.UpdateRequest) (*storage.UpdateReply, error) {
	log.Println("starting update...")

	if r == nil || r.Namespace == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad request")
	}

	nsi := snToN(r.Namespace)
	var env string
	switch r.Environment {
	case storage.Environment_Staging:
		env = environmentStaging
	case storage.Environment_Production:
		env = environmentProduction
	default:
		return nil, grpc.Errorf(codes.InvalidArgument, "bad environment provided")
	}
	err := s.sess.DB(s.db).C(env).Update(bson.M{"name": r.Namespace.Name}, nsi)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "error updating namespaces: %s", err)
	}

	return &storage.UpdateReply{Namespace: r.Namespace}, nil
}

func (s *Server) Delete(ctx context.Context, r *storage.DeleteRequest) (*storage.DeleteReply, error) {
	log.Println("starting delete...")

	if r == nil || r.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad request")
	}

	var env string
	switch r.Environment {
	case storage.Environment_Staging:
		env = environmentStaging
	case storage.Environment_Production:
		env = environmentProduction
	}

	ns, err := s.getNamespace(r.Name, env)
	if err != nil {
		return nil, err
	}

	err = s.sess.DB(s.db).C(env).Remove(bson.M{"name": r.Name})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "namespace could not be deleted: %s", err)
	}

	return &storage.DeleteReply{Namespace: ns}, nil
}

func parseNamespaces(namespaces []types.Namespace) ([]*storage.Namespace, error) {
	results := make([]*storage.Namespace, len(namespaces))
	for i, mns := range namespaces {
		ns, err := nToSN(mns)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "could not transform mongo namespace: %s", err)
		}
		results[i] = ns
	}
	return results, nil
}

func (s *Server) getNamespace(name, environment string) (*storage.Namespace, error) {
	var n types.Namespace
	err := s.sess.DB(s.db).C(environment).Find(bson.M{"name": name}).One(&n)
	if err == mgo.ErrNotFound {
		return nil, grpc.Errorf(codes.NotFound, "namespace not found: %s", name)
	} else if err != nil {
		switch err := err.(type) {
		case *mgo.QueryError:
			return nil, grpc.Errorf(codes.Internal, "query error: %s", err)
		default:
			return nil, grpc.Errorf(codes.Internal, "could not find namespace %v", name)
		}
	}

	return nToSN(n)
}

func nToSN(n types.Namespace) (*storage.Namespace, error) {
	ns := &storage.Namespace{
		Name:        n.Name,
		Labels:      n.Labels,
		Experiments: make([]*storage.Experiment, len(n.Experiments)),
	}
	for i, mexp := range n.Experiments {
		exp, err := eToSE(mexp)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "cound not transform mongo experiment: %s", err)
		}
		ns.Experiments[i] = exp
	}
	return ns, nil
}

func eToSE(e types.Experiment) (*storage.Experiment, error) {
	exp := &storage.Experiment{
		Name:   e.Name,
		Params: make([]*storage.Param, len(e.Params)),
	}
	var err error
	exp.Segments, err = decodeSegments(e.Segments)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "could not decode experiment segments: %s", err)
	}
	for i, mparam := range e.Params {
		exp.Params[i] = pToSP(mparam)
	}
	return exp, nil
}

func pToSP(p types.Param) *storage.Param {
	param := &storage.Param{
		Name: p.Name,
		Value: &storage.Value{
			Choices: p.Value.Choices,
			Weights: p.Value.Weights,
		},
	}
	return param
}

func snToN(n *storage.Namespace) types.Namespace {
	ns := types.Namespace{
		Name:        n.Name,
		Labels:      n.Labels,
		Experiments: make([]types.Experiment, len(n.Experiments)),
	}
	for i, exp := range n.Experiments {
		ns.Experiments[i] = seToE(exp)
	}
	return ns
}

func seToE(e *storage.Experiment) types.Experiment {
	exp := types.Experiment{
		Name:     e.Name,
		Segments: hex.EncodeToString(e.Segments),
		Params:   make([]types.Param, len(e.Params)),
	}
	for i, param := range e.Params {
		exp.Params[i] = spToP(param)
	}
	return exp
}

func spToP(p *storage.Param) types.Param {
	return types.Param{
		Name: p.Name,
		Value: types.Value{
			Choices: p.Value.Choices,
			Weights: p.Value.Weights,
		},
	}
}
