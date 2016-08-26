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

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"

	"github.com/foolusion/choices"
	"github.com/foolusion/choices/elwin"
	"github.com/foolusion/choices/storage/mongo"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var config = struct {
	ec              *choices.ChoicesConfig
	grpcAddr        string
	jsonAddr        string
	mongoAddr       string
	mongoDB         string
	mongoCollection string
}{
	jsonAddr:        ":8081",
	grpcAddr:        ":8080",
	mongoAddr:       "elwin-storage",
	mongoDB:         "elwin",
	mongoCollection: "test",
}

var (
	jsonRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "elwin",
		Name:      "json_requests",
		Help:      "The number of json requests recieved.",
	})
	jsonDurations = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "nordstrom",
		Subsystem: "elwin",
		Name:      "json_durations_nanoseconds",
		Help:      "json latency distributions.",
	})
	paramCounts = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "elwin",
		Name:      "param_counts",
		Help:      "Params served to users.",
	}, []string{"exp", "param", "value"})
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthzHandler)
	prometheus.MustRegister(jsonRequests)
	prometheus.MustRegister(jsonDurations)
	prometheus.MustRegister(paramCounts)
}

const (
	envJSONAddr  = "JSON_ADDRESS"
	envGRPCAddr  = "GRPC_ADDRESS"
	envMongoAddr = "MONGO_ADDRESS"
	envMongoDB   = "MONGO_DATABASE"
	envMongoColl = "MONGO_COLLECTION"
)

func main() {
	log.Println("Starting elwin...")

	// TODO: read environment variables for config
	if os.Getenv(envJSONAddr) != "" {
		config.jsonAddr = os.Getenv(envJSONAddr)
	}

	if os.Getenv(envGRPCAddr) != "" {
		config.grpcAddr = os.Getenv(envGRPCAddr)
	}

	if os.Getenv(envMongoAddr) != "" {
		config.mongoAddr = os.Getenv(envMongoAddr)
	}

	if os.Getenv(envMongoDB) != "" {
		config.mongoDB = os.Getenv(envMongoDB)
	}

	if os.Getenv(envMongoColl) != "" {
		config.mongoCollection = os.Getenv(envMongoColl)
	}

	// create elwin config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ec, err := choices.NewChoices(
		ctx,
		mongo.WithMongoStorage(config.mongoAddr, config.mongoDB, config.mongoCollection),
		choices.UpdateInterval(time.Minute),
	)
	if err != nil {
		log.Fatal(err)
	}
	config.ec = ec

	// TODO: remove when deployed
	m := ec.Storage.(*mongo.Mongo)
	m.LoadExampleData()

	http.Handle("/metrics", prometheus.Handler())

	go func() {
		config.ec.ErrChan <- http.ListenAndServe(config.jsonAddr, nil)
	}()

	go func() {
		lis, err := net.Listen("tcp", config.grpcAddr)
		if err != nil {
			config.ec.ErrChan <- fmt.Errorf("main: failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		elwin.RegisterElwinServer(grpcServer, &elwinServer{})
		config.ec.ErrChan <- grpcServer.Serve(lis)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-config.ec.ErrChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Printf("Captured %v. Exitting...", s)
			// send StatusServiceUnavailable to new requestors
			// block server from accepting new requests
			os.Exit(0)
		}
	}
}

type elwinServer struct{}

func (e *elwinServer) GetNamespaces(ctx context.Context, id *elwin.Identifier) (*elwin.Experiments, error) {
	if id == nil {
		return nil, fmt.Errorf("GetNamespaces: no Identifier recieved")
	}

	resp, err := config.ec.Namespaces(id.TeamID, id.UserID)
	if err != nil {
		return nil, fmt.Errorf("error resolving namespaces for %s, %s: %v", id.TeamID, id.UserID, err)
	}

	exp := &elwin.Experiments{
		Experiments: make(map[string]*elwin.Experiment, len(resp)),
	}

	for _, v := range resp {
		exp.Experiments[v.Name] = &elwin.Experiment{
			Params: make([]*elwin.Param, len(v.Params)),
		}

		for i, p := range v.Params {
			exp.Experiments[v.Name].Params[i] = &elwin.Param{
				Name:  p.Name,
				Value: p.Value,
			}
		}
	}
	return exp, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	r.ParseForm()
	if r.Body != nil {
		defer r.Body.Close()
	}
	resp, err := config.ec.Namespaces(r.Form.Get("teamid"), r.Form.Get("userid"))
	if err != nil {
		config.ec.ErrChan <- fmt.Errorf("rootHandler: couldn't get Namespaces: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := elwinJSON(resp, w); err != nil {
		config.ec.ErrChan <- err
		return
	}
	jsonRequests.Inc()
	jsonDurations.Observe(float64(time.Since(start)))
}

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) write(buf []byte) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write(buf)
}

var (
	jsonQuote    = []byte{'"'}
	jsonKeyEnd   = []byte{':', ' '}
	jsonComma    = []byte{','}
	jsonOpenObj  = []byte{'{'}
	jsonCloseObj = []byte{'}'}
)

func elwinJSON(er []choices.ExperimentResponse, w io.Writer) error {
	ew := &errWriter{w: w}
	ew.write(jsonOpenObj)
	ew.write(jsonQuote)
	ew.write([]byte("experiments"))
	ew.write(jsonQuote)
	ew.write(jsonKeyEnd)
	ew.write(jsonOpenObj)
	for i, exp := range er {
		ew.write(jsonQuote)
		ew.write([]byte(exp.Name))
		ew.write(jsonQuote)
		ew.write(jsonKeyEnd)
		ew.write(jsonOpenObj)
		for j, param := range exp.Params {
			ew.write(jsonQuote)
			ew.write([]byte(param.Name))
			ew.write(jsonQuote)
			ew.write(jsonKeyEnd)
			ew.write(jsonQuote)
			ew.write([]byte(param.Value))
			ew.write(jsonQuote)
			if j < len(exp.Params)-1 {
				ew.write(jsonComma)
			}
			paramCounts.With(prometheus.Labels{"exp": exp.Name, "param": param.Name, "value": param.Value}).Inc()
		}
		ew.write(jsonCloseObj)
		if i < len(er)-1 {
			ew.write(jsonComma)
		}
	}
	ew.write(jsonCloseObj)
	ew.write(jsonCloseObj)
	if ew.err != nil {
		return ew.err
	}
	return nil
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("I’ve taken on enough new tech in this project, and still need to integrate with SnowPlow. I can’t afford the time to work through the growing pains of a new code base when there is only one engineer3 backing it and they’re using completely unsupported infrastructure.\n\nLove,\nPaul McCallick"))
}
