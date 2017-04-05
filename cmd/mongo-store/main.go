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

package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	cfgListenAddr    = "listen_address"
	cfgMongoAddr     = "mongo_address"
	cfgMongoDatabase = "mongo_database"
	cfgMongoUsername = "mongo_username"
	cfgMongoPassword = "mongo_password"

	collExperiments = "experiments"
	collNamespaces  = "namespaces"
)

func main() {
	log.Println("Starting mongo-store...")
	viper.SetDefault(cfgListenAddr, ":8080")
	viper.SetDefault(cfgMongoAddr, "elwin-mongo")
	viper.SetDefault(cfgMongoDatabase, "elwin")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/mongo")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config: %v", err)
	}

	log.Printf("connecting to %s/%s", viper.GetString(cfgMongoAddr), viper.GetString(cfgMongoDatabase))

	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{viper.GetString(cfgMongoAddr)},
		Username: viper.GetString(cfgMongoUsername),
		Password: viper.GetString(cfgMongoPassword),
		Database: viper.GetString(cfgMongoDatabase),
		Timeout:  2 * time.Minute,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()
	srv := &server{
		Session: sess,
	}

	lis, err := net.Listen("tcp", viper.GetString(cfgListenAddr))
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	log.Printf("listening on %s...", viper.GetString(cfgListenAddr))

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	storage.RegisterElwinStorageServer(s, srv)
	grpc_prometheus.Register(s)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	log.Fatal(s.Serve(lis))
}

type experiment struct {
	storage.Experiment
	Id string `bson:"_id"`
}

type server struct {
	*mgo.Session
}

func (s *server) setExperiment(e experiment) error {
	_, err := s.DB(viper.GetString(cfgMongoDatabase)).
		C(collExperiments).
		UpsertId(e.Id, e)
	return err
}

func (s *server) getExperiment(id string) (*experiment, error) {
	var e experiment
	err := s.DB(viper.GetString(cfgMongoDatabase)).
		C(collExperiments).
		FindId(id).
		One(&e)
	if err != nil {
		return nil, errors.Wrap(err, "could not find experiment")
	}
	return &e, nil
}

func (s *server) List(ctx context.Context, r *storage.ListRequest) (*storage.ListReply, error) {
	if r == nil {
		return nil, errors.New("request was nil")
	}
	sel, err := labels.Parse(r.Query)
	if err != nil {
		return nil, errors.New("could not parse l")
	}
	req, selectable := sel.Requirements()
	if !selectable {
		return nil, errors.New("requirements are not selectable")
	}
	log.Println(req)
	l := bson.M{}
	for _, v := range req {
		switch v.Operator() {
		case selection.Equals, selection.DoubleEquals:
			l["experiment.labels."+v.Key()] = v.Values().List()[0]
		case selection.In:
			l["experiment.labels."+v.Key()] = bson.M{"$in": v.Values().List()}
		default:
			continue
		}
	}
	log.Println(l)
	var iter *mgo.Iter
	if len(req) == 0 {
		iter = s.DB(viper.GetString(cfgMongoDatabase)).
			C(collExperiments).
			Find(bson.M{}).
			Iter()
	} else {
		iter = s.Session.DB(viper.GetString(cfgMongoDatabase)).
			C(collExperiments).
			Find(l).
			Iter()
	}
	var exps []*storage.Experiment
	var exp experiment
	for iter.Next(&exp) {
		var a storage.Experiment
		a = exp.Experiment
		exps = append(exps, &a)
	}
	return &storage.ListReply{Experiments: exps}, nil
}

func (s *server) Get(ctx context.Context, r *storage.GetRequest) (*storage.GetReply, error) {
	if r == nil {
		return nil, errors.New("request was nil")
	}

	exp, err := s.getExperiment(r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get experiment")
	}
	return &storage.GetReply{Experiment: &exp.Experiment}, nil
}

func (s *server) Set(ctx context.Context, r *storage.SetRequest) (*storage.SetReply, error) {
	if r == nil {
		return nil, errors.New("request was nil")
	}

	if r.Experiment == nil {
		return nil, errors.New("experiment was nil")
	}

	e := experiment{
		Id:         r.Experiment.Id,
		Experiment: *r.Experiment,
	}
	if _, err := s.DB(viper.GetString(cfgMongoDatabase)).
		C(collExperiments).
		UpsertId(e.Id, e); err != nil {
		return nil, errors.Wrap(err, "could not set experiment")
	}

	return &storage.SetReply{Experiment: r.Experiment}, nil
}

func (s *server) Remove(ctx context.Context, r *storage.RemoveRequest) (*storage.RemoveReply, error) {
	if r == nil {
		return nil, errors.New("request was nil")
	}

	var exp experiment
	if err := s.DB(viper.GetString(cfgMongoDatabase)).
		C(collExperiments).
		FindId(r.Id).
		One(&exp); err != nil {
		return nil, errors.Wrap(err, "could not find record")
	}
	if err := s.DB(viper.GetString(cfgMongoDatabase)).
		C(collExperiments).
		RemoveId(r.Id); err != nil {
		return nil, errors.Wrap(err, "could not delete record")
	}
	return &storage.RemoveReply{Experiment: &exp.Experiment}, nil
}
