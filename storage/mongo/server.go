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
	"fmt"

	"google.golang.org/grpc"

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

func (s *Server) All(ctx context.Context, r *storage.AllRequest, opts ...grpc.CallOption) (*storage.AllReply, error) {
	var env string
	switch {
	case r == nil:
		env = environmentStaging
	case r.Environment == storage.Environment_Staging:
		env = environmentStaging
	case r.Environment == storage.Environment_Production:
		env = environmentProduction
	default:
		return nil, fmt.Errorf("bad environment requested")
	}

	var results []types.Namespace
	err := s.sess.DB(s.db).C(env).Find(nil).All(&results)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode data from mongo")
	}

	resp, err := parseNamespaces(results)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse namespaces")
	}

	return &storage.AllReply{
		Namespaces: resp,
	}, nil
}

func (s *Server) Create(ctx context.Context, r *storage.CreateRequest, opts ...grpc.CallOption) (*storage.CreateReply, error) {
	if r == nil || r.Namespace == nil {
		return nil, fmt.Errorf("bad request")
	}

	nsi := snToN(r.Namespace)

	var env string
	switch r.Environment {
	case storage.Environment_Staging:
		env = environmentStaging
	case storage.Environment_Production:
		env = environmentProduction
	}

	if err := s.sess.DB(s.db).C(env).Insert(nsi); err != nil {
		return nil, errors.Wrap(err, "unable to insert experiment into database")
	}
	ns, err := s.getNamespace(r.Namespace.Name, env)
	return &storage.CreateReply{Namespace: ns}, err
}

func (s *Server) Read(ctx context.Context, r *storage.ReadRequest, opts ...grpc.CallOption) (*storage.ReadReply, error) {
	if r == nil || r.Name == "" {
		return nil, fmt.Errorf("bad request")
	}
	var env string
	switch r.Environment {
	case storage.Environment_Staging:
		env = environmentStaging
	case storage.Environment_Production:
		env = environmentProduction
	}

	ns, err := s.getNamespace(r.Name, env)
	return &storage.ReadReply{Namespace: ns}, err
}

func (s *Server) Update(ctx context.Context, r *storage.UpdateRequest, opts ...grpc.CallOption) (*storage.UpdateReply, error) {
	return nil, nil
}

func (s *Server) Delete(ctx context.Context, r *storage.DeleteRequest, opts ...grpc.CallOption) (*storage.DeleteReply, error) {
	return nil, nil
}

func parseNamespaces(namespaces []types.Namespace) ([]*storage.Namespace, error) {
	results := make([]*storage.Namespace, len(namespaces))
	for i, mns := range namespaces {
		ns, err := nToSN(mns)
		if err != nil {
			return nil, errors.Wrap(err, "could not transform mongo namespace")
		}
		results[i] = ns
	}
	return results, nil
}

func (s *Server) getNamespace(name, environment string) (*storage.Namespace, error) {
	var n types.Namespace
	if err := s.sess.DB(s.db).C(environment).Find(bson.M{"name": name}).One(&n); err != nil {
		return nil, errors.Wrapf(err, "could not find namespace %v", name)
	}
	return nToSN(n)
}

func nToSN(n types.Namespace) (*storage.Namespace, error) {
	ns := &storage.Namespace{
		Name:        n.Name,
		Experiments: make([]*storage.Experiment, len(n.Experiments)),
	}
	for i, mexp := range n.Experiments {
		exp, err := eToSE(mexp)
		if err != nil {
			return nil, errors.Wrap(err, "cound not transform mongo experiment")
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
		return nil, errors.Wrap(err, "could not decode experiment segments")
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
