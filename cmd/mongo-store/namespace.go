package main

import (
	"context"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type namespace struct {
	Name string `bson:"_id"`
	*storage.Namespace
}

func (s *server) Namespace(ctx context.Context, name string) (*storage.Namespace, error) {
	var ns namespace
	if err := s.DB(s.db).C(collNamespaces).FindId(name).One(&ns); err != nil {
		return nil, errors.Wrap(err, "could not find namespace")
	}
	return ns.Namespace, nil
}

func (s *server) SetNamespace(ctx context.Context, n *storage.Namespace) error {
	ns := namespace{
		Name:      n.Name,
		Namespace: n,
	}
	if _, err := s.DB(s.db).C(collNamespaces).UpsertId(ns.Name, ns); err != nil {
		return errors.Wrap(err, "could not update namespace")
	}
	return nil
}

func (s *server) AllNamespaces(ctx context.Context) ([]*storage.Namespace, error) {
	var namespaces []namespace
	if err := s.DB(s.db).C(collNamespaces).Find(bson.M{}).All(&namespaces); err != nil {
		return nil, errors.Wrap(err, "could not find all namespaces")
	}
	n := make([]*storage.Namespace, len(namespaces))
	for i := range n {
		n[i] = namespaces[i].Namespace
	}
	return n, nil
}
