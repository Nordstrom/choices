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
	cfgListenAddr      = "listen_address"
	cfgMongoAddr       = "mongo_address"
	cfgMongoDatabase   = "mongo_database"
	cfgMongoCollection = "mongo_collection"
	cfgMongoUsername   = "mongo_username"
	cfgMongoPassword   = "mongo_password"
)

func main() {
	log.Println("Starting mongo-store...")
	viper.SetDefault(cfgListenAddr, ":8080")
	viper.SetDefault(cfgMongoAddr, "elwin-mongo")
	viper.SetDefault(cfgMongoDatabase, "elwin")
	viper.SetDefault(cfgMongoCollection, "dev")

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
		session: sess,
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
	ID string `bson:"_id"`
	storage.Experiment
}

type server struct {
	session *mgo.Session
	c       *mgo.Collection
}

func getCollection(s *server) *mgo.Collection {
	return s.session.DB(viper.GetString(cfgMongoDatabase)).C(viper.GetString(cfgMongoCollection))
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
	c := getCollection(s)

	var iter *mgo.Iter
	if len(req) == 0 {
		iter = c.Find(bson.M{}).Iter()
	} else {
		iter = c.Find(l).Iter()
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

	c := getCollection(s)
	var exp experiment
	if err := c.FindId(r.Id).One(&exp); err != nil {
		return nil, errors.Wrap(err, "could not find record")
	}
	return &storage.GetReply{Experiment: &exp.Experiment}, nil
}

func (s *server) Set(ctx context.Context, r *storage.SetRequest) (*storage.SetReply, error) {
	if r == nil {
		return nil, errors.New("request was nil")
	}

	c := getCollection(s)

	if r.Experiment == nil {
		return nil, errors.New("experiment was nil")
	}

	e := experiment{
		ID:         r.Experiment.Id,
		Experiment: *r.Experiment,
	}
	if _, err := c.UpsertId(e.ID, e); err != nil {
		return nil, errors.Wrap(err, "could not set experiment")
	}

	return &storage.SetReply{Experiment: r.Experiment}, nil
}

func (s *server) Remove(ctx context.Context, r *storage.RemoveRequest) (*storage.RemoveReply, error) {
	if r == nil {
		return nil, errors.New("request was nil")
	}

	c := getCollection(s)
	var exp experiment
	if err := c.FindId(r.Id).One(&exp); err != nil {
		return nil, errors.Wrap(err, "could not find record")
	}
	if err := c.RemoveId(r.Id); err != nil {
		return nil, errors.Wrap(err, "could not delete record")
	}
	return &storage.RemoveReply{Experiment: &exp.Experiment}, nil
}
