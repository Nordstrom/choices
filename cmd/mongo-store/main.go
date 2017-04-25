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
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Nordstrom/choices"
	"github.com/foolusion/elwinprotos/storage"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	oldctx "golang.org/x/net/context"
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

var (
	errNilRequest = errors.New("request was nil")
)

type experiment struct {
	*storage.Experiment
	ID string `bson:"_id"`
}

type server struct {
	*mgo.Session
	db string
}

func (s *server) SetExperiment(ctx context.Context, e *storage.Experiment) error {
	s.Refresh()
	if e == nil {
		return errors.New("experiment is nil")
	}
	exp := experiment{
		Experiment: e,
		ID:         e.Id,
	}
	_, err := s.DB(s.db).C(collExperiments).UpsertId(e.Id, exp)
	return err
}

func (s *server) Experiment(ctx context.Context, id string) (*storage.Experiment, error) {
	s.Refresh()
	var e experiment
	if err := s.DB(s.db).C(collExperiments).FindId(id).One(&e); err != nil {
		return nil, errors.Wrap(err, "could not find experiment")
	}
	return e.Experiment, nil
}

func (s *server) AllExperiments(ctx context.Context) ([]*storage.Experiment, error) {
	s.Refresh()
	resp, err := s.List(ctx, &storage.ListRequest{})
	if err != nil {
		return nil, nil
	}
	return resp.Experiments, nil
}

func (s *server) New(ctx oldctx.Context, r *storage.NewRequest) (*storage.NewReply, error) {
	if r == nil {
		return nil, errNilRequest
	}
	if err := choices.CreateExperiment(ctx, s, r.Experiment, r.Namespace, int(r.NSegments), int(r.ESegments)); err != nil {
		return nil, errors.Wrap(err, "could not create experiment")
	}
	return &storage.NewReply{}, nil
}

func (s *server) List(ctx oldctx.Context, r *storage.ListRequest) (*storage.ListReply, error) {
	s.Refresh()
	if r == nil {
		return nil, errNilRequest
	}
	sel, err := labels.Parse(r.Query)
	if err != nil {
		return nil, errors.New("could not parse labels")
	}
	req, selectable := sel.Requirements()
	if !selectable {
		return nil, errors.New("requirements are not selectable")
	}
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
	var iter *mgo.Iter
	if len(req) == 0 {
		iter = s.DB(s.db).C(collExperiments).Find(bson.M{}).Iter()
	} else {
		iter = s.Session.DB(s.db).C(collExperiments).Find(l).Iter()
	}
	var exps []*storage.Experiment
	var exp experiment
	for iter.Next(&exp) {
		exps = append(exps, exp.Experiment)
	}
	return &storage.ListReply{Experiments: exps}, nil
}

func (s *server) Get(ctx oldctx.Context, r *storage.GetRequest) (*storage.GetReply, error) {
	s.Refresh()
	if r == nil {
		return nil, errNilRequest
	}

	exp, err := s.Experiment(ctx, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get experiment")
	}
	return &storage.GetReply{Experiment: exp}, nil
}

func (s *server) Set(ctx oldctx.Context, r *storage.SetRequest) (*storage.SetReply, error) {
	s.Refresh()
	if r == nil {
		return nil, errNilRequest
	}

	if r.Experiment == nil {
		return nil, errors.New("experiment is nil")
	}

	var backup *storage.Experiment
	if b, err := s.Experiment(ctx, r.Experiment.Id); err != nil {
		// probably just not found
		log.Println(err)
	} else {
		backup = b
	}

	if err := s.SetExperiment(ctx, r.Experiment); err != nil {
		return nil, errors.Wrap(err, "could not set experiment")
	}

	err := choices.AutoFix(ctx, s)
	if err != nil && backup != nil {
		// try to restore backup
		if err2 := s.SetExperiment(ctx, backup); err2 != nil {
			return nil, errors.Wrap(errors.Wrap(err, err2.Error()), "restore previous version failed")
		}
		return nil, errors.Wrap(err, "bad experiment: restored previous version")
	} else if err != nil {
		// just delete the bad experiment
		if _, err := s.Remove(ctx, &storage.RemoveRequest{Id: r.Experiment.Id}); err != nil {
			return nil, errors.Wrap(err, "could not delete bad experiment")
		}
		return nil, errors.Wrap(err, "bad experiment: removed experiment")
	}

	return &storage.SetReply{Experiment: r.Experiment}, nil
}

func (s *server) Remove(ctx oldctx.Context, r *storage.RemoveRequest) (*storage.RemoveReply, error) {
	s.Refresh()
	if r == nil {
		return nil, errNilRequest
	}

	exp, err := s.Experiment(ctx, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get experiment")
	}

	if err := s.DB(s.db).C(collExperiments).RemoveId(r.Id); err != nil {
		return nil, errors.Wrap(err, "could not delete record")
	}
	return &storage.RemoveReply{Experiment: exp}, nil
}

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
		db:      viper.GetString(cfgMongoDatabase),
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
