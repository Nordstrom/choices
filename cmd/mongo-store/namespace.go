package main

import (
	"context"

	"github.com/Nordstrom/choices"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type namespace struct {
	Name string `bson:"_id"`
	choices.Namespace
}

func (s *server) Namespace(ctx context.Context, name string) (*choices.Namespace, error) {
	var ns namespace
	if err := s.DB(s.db).C(collNamespaces).FindId(name).One(&ns); err != nil {
		return nil, errors.Wrap(err, "could not find namespace")
	}
	return &ns.Namespace, nil
}

func (s *server) SetNamespace(ctx context.Context, n *choices.Namespace) error {
	ns := namespace{
		Name:      n.Name,
		Namespace: *n,
	}
	if _, err := s.DB(s.db).C(collNamespaces).UpsertId(ns.Name, ns); err != nil {
		return errors.Wrap(err, "could not update namespace")
	}
	return nil
}

func (s *server) AllNamespaces(ctx context.Context) ([]*choices.Namespace, error) {
	var namespaces []namespace
	if err := s.DB(s.db).C(collNamespaces).Find(bson.M{}).All(namespaces); err != nil {
		return nil, errors.Wrap(err, "could not find all namespaces")
	}
	n := make([]*choices.Namespace, len(namespaces))
	for i := range n {
		n[i] = &namespaces[i].Namespace
	}
	return n, nil
}
