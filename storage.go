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

package choices

import (
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	storage "github.com/foolusion/choices/elwinstorage"
	"github.com/pkg/errors"
)

const (
	StorageEnvironmentBad = iota
	StorageEnvironmentDev
	StorageEnvironmentProd
)

var ErrBadStorageEnvironment = errors.New("bad storage environment")

func WithStorageConfig(addr string, env int) ConfigOpt {
	return func(c *Config) error {
		cc, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return errors.Wrap(err, "could not dial storage service")
		}
		if env == StorageEnvironmentBad {
			return ErrBadStorageEnvironment
		}
		c.Storage = NewNamespaceStore(cc, env)
		return nil
	}
}

type NamespaceStore struct {
	mu    sync.RWMutex
	el    storage.ElwinStorageClient
	env   int
	cache []Namespace
}

func NewNamespaceStore(cc *grpc.ClientConn, env int) *NamespaceStore {
	return &NamespaceStore{
		el:  storage.NewElwinStorageClient(cc),
		env: env,
	}
}

func (n *NamespaceStore) Read() []Namespace {
	out := make([]Namespace, len(n.cache))
	n.mu.RLock()
	copy(out, n.cache)
	n.mu.RUnlock()
	return out
}

func (n *NamespaceStore) Update() error {
	var req *storage.AllRequest
	switch n.env {
	case StorageEnvironmentDev:
		req = &storage.AllRequest{
			Environment: storage.Environment_Staging,
		}
	case StorageEnvironmentProd:
		req = &storage.AllRequest{
			Environment: storage.Environment_Production,
		}
	default:
		return ErrBadStorageEnvironment
	}
	ar, err := n.el.All(context.TODO(), req)
	if err != nil {
		return errors.Wrap(err, "error requesting All from storage")
	}

	cache := make([]Namespace, len(ar.GetNamespaces()))
	for i, ns := range ar.GetNamespaces() {
		var err error
		cache[i], err = FromNamespace(ns)
		if err != nil {
			return errors.Wrap(err, "could not parse Namespace")
		}
	}
	n.mu.Lock()
	n.cache = cache
	n.mu.Unlock()
	return nil
}

func FromNamespace(s *storage.Namespace) (Namespace, error) {
	ns := NewNamespace(s.Name, s.Labels)
	for _, e := range s.Experiments {
		err := ns.AddExperiment(FromExperiment(e))
		if err != nil {
			return Namespace{}, errors.Wrap(err, "could not remove add experiment")
		}
	}
	return *ns, nil
}

func FromExperiment(s *storage.Experiment) Experiment {
	exp := Experiment{
		Name:   s.Name,
		Params: make([]Param, len(s.Params)),
	}
	copy(exp.Segments[:], s.Segments[:16])

	for i, p := range s.Params {
		exp.Params[i] = FromParam(p)
	}

	return exp
}

func FromParam(s *storage.Param) Param {
	par := Param{
		Name: s.Name,
	}
	switch {
	case len(s.Value.Weights) == 0:
		par.Value = &Uniform{
			Choices: s.Value.Choices,
		}
	case len(s.Value.Weights) == len(s.Value.Choices):
		par.Value = &Weighted{
			Choices: s.Value.Choices,
			Weights: s.Value.Weights,
		}
	}
	return par
}

// TeamNamespaces filters the namespaces from storage based on teamID.
func TeamNamespaces(s NamespaceStore, teamID string) []Namespace {
	allNamespaces := s.Read()
	teamNamespaces := make([]Namespace, 0, len(allNamespaces))
	for _, n := range allNamespaces {
		for _, t := range n.TeamID {
			if t == teamID {
				teamNamespaces = append(teamNamespaces, n)
			}
		}
	}
	return teamNamespaces
}
