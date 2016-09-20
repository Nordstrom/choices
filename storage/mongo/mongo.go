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
	"sync"

	"github.com/foolusion/choices"
	"github.com/foolusion/choices/storage/mongo/internal/types"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo implements the storage interface.
type Mongo struct {
	namespaces    []choices.Namespace
	sess          *mgo.Session
	url, db, coll string
	mu            sync.RWMutex
}

// WithMongoStorage is a helper func to create a mongo connection and add the data to the choices.Config
func WithMongoStorage(url, db, collection string) func(*choices.Config) error {
	return func(ec *choices.Config) error {
		m := &Mongo{url: url, db: db, coll: collection}
		sess, err := mgo.Dial(url)
		if err != nil {
			return err
		}
		m.sess = sess
		ec.Storage = m
		return nil
	}
}

// Read returns the current namespaces stored in the mongo object.
func (m *Mongo) Read() []choices.Namespace {
	m.mu.RLock()
	ns := m.namespaces
	m.mu.RUnlock()
	return ns
}

// Ready returns whether the mongo database is available.
func (m *Mongo) Ready() error {
	return m.sess.Ping()
}

// Update updates the data in the mongo object.
func (m *Mongo) Update() error {
	c := m.sess.DB(m.db).C(m.coll)
	iter := c.Find(bson.M{}).Iter()
	var mongoNamespaces []types.Namespace
	err := iter.All(&mongoNamespaces)
	if err != nil {
		return err
	}

	namespaces := make([]choices.Namespace, len(mongoNamespaces))
	for i, n := range mongoNamespaces {
		namespaces[i], err = parseNamespace(n)
		if err != nil {
			return err
		}
	}

	m.mu.Lock()
	m.namespaces = namespaces
	m.mu.Unlock()
	return nil
}

// NamespaceToChoicesNamespace converts the data read from mongo into a proper choices data structure.
func NamespaceToChoicesNamespace(n types.Namespace) (choices.Namespace, error) {
	return parseNamespace(n)
}

func decodeSegments(seg string) ([]byte, error) {
	segBytes, err := hex.DecodeString(seg)
	if err != nil {
		return nil, err
	}
	return segBytes[:16], nil
}

func parseNamespace(n types.Namespace) (choices.Namespace, error) {
	namespace := choices.Namespace{
		Name:        n.Name,
		TeamID:      n.TeamID,
		Experiments: make([]choices.Experiment, len(n.Experiments)),
	}
	nss, err := decodeSegments(n.Segments)
	if err != nil {
		return choices.Namespace{}, err
	}
	var buf [16]byte
	copy(buf[:], nss[:16])
	namespace.Segments = buf
	for i, e := range n.Experiments {
		namespace.Experiments[i], err = parseExperiment(e)
		if err != nil {
			return choices.Namespace{}, err
		}
	}
	return namespace, nil
}

func parseExperiment(e types.Experiment) (choices.Experiment, error) {
	experiment := choices.Experiment{
		Name:   e.Name,
		Params: make([]choices.Param, len(e.Params)),
	}
	ess, err := decodeSegments(e.Segments)
	if err != nil {
		return choices.Experiment{}, err
	}
	var buf [16]byte
	copy(buf[:], ess[:16])
	experiment.Segments = buf

	for i, p := range e.Params {
		experiment.Params[i] = parseParam(p)
	}
	return experiment, nil
}

func parseParam(p types.Param) choices.Param {
	var param choices.Param
	param = choices.Param{Name: p.Name}
	switch p.Type {
	case choices.ValueTypeUniform:
		var uniform choices.Uniform
		p.Value.Unmarshal(&uniform)
		param.Value = &uniform
	case choices.ValueTypeWeighted:
		var weighted choices.Weighted
		p.Value.Unmarshal(&weighted)
		param.Value = &weighted
	}
	return param
}

// QueryAll querys the namespaces using the given query and returns all matches.
func QueryAll(c *mgo.Collection, query interface{}) ([]choices.Namespace, error) {
	iter := c.Find(query).Iter()
	var mongoNamespaces []types.Namespace
	err := iter.All(&mongoNamespaces)
	if err != nil {
		return nil, err
	}

	namespaces := make([]choices.Namespace, len(mongoNamespaces))
	for i, n := range mongoNamespaces {
		namespaces[i], err = parseNamespace(n)
		if err != nil {
			return nil, err
		}
	}

	return namespaces, nil
}

// QueryOne querys the namespace using the given query and returns the first match.
func QueryOne(c *mgo.Collection, query interface{}) (choices.Namespace, error) {
	var mongoNamespace types.Namespace
	if err := c.Find(query).One(&mongoNamespace); err != nil {
		return choices.Namespace{}, err
	}
	return parseNamespace(mongoNamespace)
}

// Upsert inserts a namespace into the database if it does not exist or updates the namespace if it does exist.
func Upsert(c *mgo.Collection, name string, namespace choices.Namespace) error {
	nsi := types.NamespaceInput{
		Name:        namespace.Name,
		TeamID:      namespace.TeamID,
		Segments:    hex.EncodeToString(namespace.Segments[:]),
		Experiments: make([]types.ExperimentInput, len(namespace.Experiments)),
	}
	for i, exp := range namespace.Experiments {
		nsi.Experiments[i] = types.ExperimentInput{
			Name:     exp.Name,
			Segments: hex.EncodeToString(exp.Segments[:]),
			Params:   make([]types.ParamInput, len(exp.Params)),
		}
		for j, param := range exp.Params {
			nsi.Experiments[i].Params[j] = types.ParamInput{
				Name:  param.Name,
				Value: param.Value,
			}
			switch param.Value.(type) {
			case *choices.Uniform:
				nsi.Experiments[i].Params[j].Type = choices.ValueTypeUniform
			case *choices.Weighted:
				nsi.Experiments[i].Params[j].Type = choices.ValueTypeWeighted
			default:
				return fmt.Errorf("bad param type")
			}
		}
	}
	if _, err := c.Upsert(bson.M{"name": name}, nsi); err != nil {
		return err
	}
	return nil
}
