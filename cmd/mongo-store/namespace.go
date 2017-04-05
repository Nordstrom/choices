package main

import (
	"math/rand"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"k8s.io/contrib/compare/Godeps/_workspace/src/github.com/opencontainers/runc/libcontainer/utils"
)

type namespace struct {
	Name        string `bson:"_id"`
	NumSegments int
	Segments    []byte
}

func (n *namespace) sampleSegments(num int) []byte {
	if num <= 0 {
		return make([]byte, len(n.Segments))
	}
	available := n.availableSegments()
	p := rand.Perm(len(available))
	out := make([]byte, len(n.Segments))
	if num > len(available) {
		num = len(available)
	}
	for i := 0; i < num; i++ {
		set(out, available[p[i]])
	}
	return out
}

func set(b []byte, index int) {
	i, pos := index/8, uint8(index%8)
	b[i] |= 1 << pos
}

func clear(b []byte, index int) {
	i, pos := index/8, uint8(index%8)
	b[i] &= ^(1 << pos)
}

func (n *namespace) availableSegments() []int {
	out := make([]int, 0, n.NumSegments)
	for i := range n.Segments {
		for shift := uint8(0); shift < 8; shift++ {
			if i*8+int(shift) > int(n.NumSegments) {
				break
			}
			if n.Segments[i]&(1<<shift) != 1<<shift {
				out = append(out, i*8+int(shift))
			}
		}
	}
	return out
}

func (s *server) getNamespace(name string) (*namespace, error) {
	var ns namespace
	if err := s.DB(viper.GetString(cfgMongoDatabase)).C(collNamespaces).FindId(name).One(&ns); err != nil {
		return nil, errors.Wrap(err, "could not find namespace")
	}
	return &ns, nil
}

func (s *server) createNamespace(name string, numSegments int) (*namespace, error) {
	if numSegments <= 0 {
		return nil, errors.New("must supply a positive number of segments")
	}
	numBytes := numSegments / 8
	if numSegments%8 != 0 {
		numBytes++
	}
	return &namespace{
		Name:        name,
		NumSegments: numSegments,
		Segments:    make([]byte, numBytes),
	}, nil
}

func (s *server) createExperiment(name, namespace string, nsSegments, expSegments int, labels map[string]string, params []*storage.Param) error {
	ns, err := s.getNamespace(namespace)
	if err != nil {
		name, err := utils.GenerateRandomName("", 7)
		if err != nil {
			return errors.Wrap(err, "could not genereate new namespace name")
		}
		if ns, err = s.createNamespace(name, nsSegments); err != nil {
			return errors.Wrap(err, "could not create new namespace")
		}
	}

	seg := ns.sampleSegments(expSegments)

	if _, err := s.DB(viper.GetString(cfgMongoDatabase)).C(collNamespaces).UpsertId(ns.Name, *ns); err != nil {
		return errors.Wrap(err, "could not update namespace")
	}

	id, err := utils.GenerateRandomName(name, 32)
	if err != nil {
		return errors.Wrap(err, "could not generate experiment id")
	}

	exp := experiment{
		Id: id,
		Experiment: storage.Experiment{
			Id:        id,
			Name:      name,
			Namespace: ns.Name,
			Labels:    labels,
			Params:    params,
			Segments:  seg,
		},
	}
	if _, err := s.DB(viper.GetString(cfgMongoDatabase)).C(collExperiments).UpsertId(id, exp); err != nil {
		return errors.Wrap(err, "could not update experiment")
	}
	return nil
}
