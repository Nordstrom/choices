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
	"context"
	"crypto/rand"
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"github.com/foolusion/choices"
	storage "github.com/foolusion/choices/elwinstorage"
	"github.com/foolusion/choices/storage/mongo/internal/types"
	"github.com/pkg/errors"
)

const (
	environmentStaging = "staging"
	environmentProd    = "production"
)

type server struct {
	DB *mgo.Database
}

func (s *server) All(ctx context.Context, r *storage.AllRequest) (*storage.NamespacesReply, error) {
	var env string
	switch r.Environment {
	case storage.Environment_BAD_ENVIRONMENT:
		env = environmentStaging
	case storage.Environment_PRODUCTION:
		env = environmentProd
	default:
		return nil, fmt.Errorf("bad environment requested")
	}

	var results []types.Namespace
	err := s.DB.C(env).Find(nil).All(&results)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode data from mongo")
	}

	resp, err := parseNamespaces(results)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse namespaces")
	}

	return &storage.NamespacesReply{
		Namespaces: resp,
	}, nil
}

func (s *server) CreateExperiment(ctx context.Context, r *storage.CreateExperimentRequest) (*storage.Namespace, error) {
	if r == nil {
		return nil, fmt.Errorf("bad request")
	}

	var nsi types.NamespaceInput
	if r.Namespace == "" {
		nsi = types.NamespaceInput{
			Name:     randomName(6),
			Segments: allSegments,
		}
	}
	err := s.DB.C(environmentStaging).Insert(nsi)
	return nil, err
}

func (s *server) DeleteExperiment(ctx context.Context, r *storage.DeleteExperimentRequest) (*storage.Namespace, error) {
	return nil, nil
}

func (s *server) PublishExperiment(ctx context.Context, r *storage.PublishExperimentRequest) (*storage.Namespace, error) {
	return nil, nil
}

func (s *server) UnpublishExperiment(ctx context.Context, r *storage.UnpublishExperimentRequest) (*storage.Namespace, error) {
	return nil, nil
}

func parseNamespaces(namespaces []types.Namespace) ([]*storage.Namespace, error) {
	results := make([]*storage.Namespace, len(namespaces))
	for i, mns := range namespaces {
		ns, err := nToN(mns)
		if err != nil {
			return nil, errors.Wrap(err, "could not transform mongo namespace")
		}
		results[i] = ns
	}
	return results, nil
}

func nToN(n types.Namespace) (*storage.Namespace, error) {
	ns := &storage.Namespace{
		Name:        n.Name,
		Experiments: make([]*storage.Experiment, len(n.Experiments)),
	}
	var err error
	ns.Segments, err = decodeSegments(n.Segments)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode namespace segments")
	}
	for j, mexp := range n.Experiments {
		exp, err := eToE(mexp)
		if err != nil {
			return nil, errors.Wrap(err, "cound not transform mongo experiment")
		}
		ns.Experiments[j] = exp
	}
	return ns, nil
}

func eToE(e types.Experiment) (*storage.Experiment, error) {
	exp := &storage.Experiment{
		Name:   e.Name,
		Params: make([]*storage.Experiment_Param, len(e.Params)),
	}
	var err error
	exp.Segments, err = decodeSegments(e.Segments)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode experiment segments")
	}
	for k, mparam := range e.Params {
		exp.Params[k] = pToP(mparam)
	}
	return exp, nil
}

func pToP(p types.Param) *storage.Experiment_Param {
	param := &storage.Experiment_Param{
		Name: p.Name,
	}
	switch p.Type {
	case choices.ValueTypeUniform:
		var u choices.Uniform
		p.Value.Unmarshal(&u)
		val := &storage.Experiment_Param_Value{
			ValueType: storage.Experiment_Param_Value_UNIFORM,
			Choices:   u.Choices,
		}
		param.Value = val
	case choices.ValueTypeWeighted:
		var w choices.Weighted
		p.Value.Unmarshal(&w)
		val := &storage.Experiment_Param_Value{
			ValueType: storage.Experiment_Param_Value_WEIGHTED,
			Choices:   w.Choices,
			Weights:   w.Weights,
		}
		param.Value = val
	}
	return param
}

var alphaNum = []byte{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func randomName(n int) string {
	buf := make([]byte, n)
	rand.Read(buf)
	for i, v := range buf {
		buf[i] = alphaNum[int(v)%len(alphaNum)]
	}
	return string(buf)
}
