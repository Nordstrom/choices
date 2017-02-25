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

package choices

import (
	"sync"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/pkg/errors"
)

// constants for storage enviroments. So far only support a staging and
// production.
const (
	StorageEnvironmentBad = iota
	StorageEnvironmentDev
	StorageEnvironmentProd
)

// ErrBadStorageEnvironment is an error for when the storage environment is not
// set correctly.
var ErrBadStorageEnvironment = errors.New("bad storage environment")

// WithStorageConfig is where you set the address and environment you'd like to
// point. This is used as a ConfigOpt in NewChoices.
func WithStorageConfig(addr string, env int, updateInterval time.Duration) ConfigOpt {
	return func(c *Config) error {
		cc, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(updateInterval))
		if err != nil {
			return errors.Wrap(err, "could not dial storage service")
		}
		c.clientConn = cc
		if env == StorageEnvironmentBad {
			return ErrBadStorageEnvironment
		}
		c.storage = newNamespaceStore(cc, env)
		return nil
	}
}

// namespaceStore is the in memory copy of the storage. el is the
// storage.ElwinStorageClient used to get the data out of storage.
type namespaceStore struct {
	mu            sync.RWMutex
	el            storage.ElwinStorageClient
	env           int
	cache         []Namespace
	failedUpdates int
	isHealthy     error
}

// newNamespaceStore creates a new in memory store for the data and client to
// use to update the in memory store.
func newNamespaceStore(cc *grpc.ClientConn, env int) *namespaceStore {
	return &namespaceStore{
		el:  storage.NewElwinStorageClient(cc),
		env: env,
	}
}

// read returns the current list of Namespace that are in memory.
func (n *namespaceStore) read() []Namespace {
	out := make([]Namespace, len(n.cache))
	n.mu.RLock()
	copy(out, n.cache)
	n.mu.RUnlock()
	return out
}

// update requests the data from storage server and updates the in memory copy
// with the lastest data. It returns wether or not the update was successful.
func (n *namespaceStore) update() error {
	var req *storage.AllRequest
	switch n.env {
	case StorageEnvironmentDev:
		req = &storage.AllRequest{
			Environment: storage.Staging,
		}
	case StorageEnvironmentProd:
		req = &storage.AllRequest{
			Environment: storage.Production,
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

// FromNamespace converts a *storage.Namespace into a Namespace.
func FromNamespace(s *storage.Namespace) (Namespace, error) {
	ns := NewNamespace(s.Name)
	for _, e := range s.Experiments {
		err := ns.AddExperiment(FromExperiment(e))
		if err != nil {
			return Namespace{}, errors.Wrap(err, "could not remove add experiment")
		}
	}
	return *ns, nil
}

// FromExperiment converts a *storage.Experiment into an Experiment
func FromExperiment(s *storage.Experiment) Experiment {
	exp := Experiment{
		Name:   s.Name,
		Params: make([]Param, len(s.Params)),
		Labels: s.Labels,
	}
	copy(exp.Segments[:], s.Segments[:16])

	for i, p := range s.Params {
		exp.Params[i] = FromParam(p)
	}

	return exp
}

// FromParam converts a *storage.Param into a Param
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
func TeamNamespaces(s *namespaceStore, teamID string) []Namespace {
	allNamespaces := s.read()
	teamNamespaces := make([]Namespace, 0, len(allNamespaces))
	for _, n := range allNamespaces {
		for _, e := range n.Experiments {
			if teamID == e.Labels["team"] {
				teamNamespaces = append(teamNamespaces, n)
				break
			}
		}
	}
	return teamNamespaces
}
