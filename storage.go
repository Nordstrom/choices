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
	"log"
	"sync"
	"time"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/labels"
)

// ErrBadStorageEnvironment is an error for when the storage environment is not
// set correctly.
var ErrBadStorageEnvironment = errors.New("bad storage environment")

// WithStorageConfig is where you set the address and environment you'd like to
// point. This is used as a ConfigOpt in NewChoices.
func WithStorageConfig(addr string, updateInterval time.Duration) ConfigOpt {
	return func(c *Config) error {
		cc, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(updateInterval))
		if err != nil {
			return errors.Wrap(err, "could not dial storage service")
		}

		c.storage = newExperimentStore(cc)
		return nil
	}
}

type erPool struct {
	mu   sync.RWMutex
	pool sync.Pool
}

func (e *erPool) get() []ExperimentResponse {
	e.mu.RLock()
	defer e.mu.RUnlock()

	out, ok := e.pool.Get().([]ExperimentResponse)
	if out == nil {
		return make([]ExperimentResponse, 0, 100)
	} else if !ok {
		log.Fatal("pool had unknown object")
	}
	return out
}

func (e *erPool) put(exp []ExperimentResponse) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.pool.Put(exp)
}

func (e *erPool) new(f func() interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.pool.New = f
}

// experimentStore is the in memory copy of the storage. el is the
// storage.ElwinStorageClient used to get the data out of storage.
type experimentStore struct {
	el storage.ElwinStorageClient

	mu            sync.RWMutex
	cache         []*Experiment
	failedUpdates int
	pool          erPool
}

// newExperimentStore creates a new in memory store for the data and client to
// use to update the in memory store.
func newExperimentStore(cc *grpc.ClientConn) *experimentStore {
	return &experimentStore{
		el: storage.NewElwinStorageClient(cc),
	}
}

// read returns the current list of Experiment that are in memory.
func (n *experimentStore) read() []*Experiment {
	out := make([]*Experiment, len(n.cache))
	n.mu.RLock()
	copy(out, n.cache)
	n.mu.RUnlock()
	return out
}

// update requests the data from storage server and updates the in memory copy
// with the latest data. It returns wether or not the update was successful.
func (n *experimentStore) update() error {
	ar, err := n.el.List(context.TODO(), &storage.ListRequest{})
	if err != nil {
		log.Print("### ERROR in storage.go -> line 79", &storage.ListRequest{})
		return errors.Wrap(err, "error requesting All from storage")
	}

	cache := make([]*Experiment, len(ar.Experiments))
	for i, exp := range ar.Experiments {
		cache[i] = FromExperiment(exp)
	}
	n.mu.Lock()
	n.cache = cache
	n.mu.Unlock()

	go n.pool.new(func() interface{} { return make([]ExperimentResponse, 0, len(n.cache)) })

	return nil
}

func FromNamespace(s *storage.Namespace) *Namespace {
	return &Namespace{
		Name:        s.Name,
		NumSegments: int(s.NumSegments),
		Segments:    FromSegments(s.Segments),
	}
}

func FromSegments(s *storage.Segments) *segments {
	if s == nil {
		return &segments{b: make([]byte, defaultNumSegments/8), len: defaultNumSegments}
	}
	return &segments{b: s.B, len: int(s.Len)}
}

// FromExperiment converts a *storage.Experiment into an Experiment
func FromExperiment(s *storage.Experiment) *Experiment {
	exp := &Experiment{
		ID:        s.Id,
		Name:      s.Name,
		Namespace: s.Namespace,
		Params:    make([]Param, len(s.Params)),
		Labels:    s.Labels,
		Segments:  FromSegments(s.Segments), // TODO: what do we do if this is nil?
	}

	for i, p := range s.Params {
		exp.Params[i] = fromParam(p)
	}

	return exp
}

// FromParam converts a *storage.Param into a Param
func fromParam(s *storage.Param) Param {
	par := Param{
		Name: s.Name,
	}
	switch {
	case len(s.Value.Weights) == 0:
		par.Choices = &Uniform{
			Choices: s.Value.Choices,
		}
	case len(s.Value.Weights) == len(s.Value.Choices):
		choices := make([]weightedChoice, len(s.Value.Choices))
		for i := range choices {
			choices[i].name = s.Value.Choices[i]
			choices[i].weight = s.Value.Weights[i]
		}
		par.Choices = &Weighted{
			Choices: choices,
		}
	}
	return par
}

// teamNamespaces filters the namespaces from storage based on teamID.
func teamNamespaces(s *experimentStore, selector labels.Selector) []*Experiment {
	experiments := s.read()
	filtered := make([]*Experiment, 0, len(experiments))
	for _, e := range experiments {
		if selector.Matches(e.Labels) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
