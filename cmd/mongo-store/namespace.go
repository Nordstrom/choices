package main

import (
	"github.com/Nordstrom/choices"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type namespace struct {
	Name string `bson:"_id"`
	choices.Namespace
}

func (s *server) Namespace(name string) (*choices.Namespace, error) {
	var ns namespace
	if err := s.DB(viper.GetString(cfgMongoDatabase)).C(collNamespaces).FindId(name).One(&ns); err != nil {
		return nil, errors.Wrap(err, "could not find namespace")
	}
	return &ns.Namespace, nil
}

func (s *server) SetNamespace(n *choices.Namespace) error {
	ns := namespace{
		Name:      n.Name,
		Namespace: *n,
	}
	if _, err := s.DB(viper.GetString(cfgMongoDatabase)).C(collNamespaces).UpsertId(ns.Name, ns); err != nil {
		return errors.Wrap(err, "could not update namespace")
	}
	return nil
}
