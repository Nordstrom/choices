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
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var config = struct {
	ec              *choices.Config
	grpcAddr        string
	jsonAddr        string
	mongoAddr       string
	mongoDB         string
	mongoCollection string
	readiness       struct {
		storage      bool
		grpcServer   bool
		errorHandler bool
		httpServer   bool
	}
}{
	jsonAddr:        ":8081",
	grpcAddr:        ":8080",
	mongoAddr:       "elwin-storage:80",
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
	http.HandleFunc("/gen", genHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readiness", readinessHandler)
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

func env(dst *string, src string) {
	if os.Getenv(src) != "" {
		*dst = os.Getenv(src)
		log.Printf("Set %s to %s", src, *dst)
	}
}

func main() {
	log.Println("Starting elwin...")

	env(&config.jsonAddr, envJSONAddr)
	env(&config.grpcAddr, envGRPCAddr)
	env(&config.mongoAddr, envMongoAddr)
	env(&config.mongoDB, envMongoDB)
	env(&config.mongoCollection, envMongoColl)

	var storageEnv int
	switch config.mongoCollection {
	case "staging", "dev", "test":
		storageEnv = choices.StorageEnvironmentDev
	case "production", "prod":
		storageEnv = choices.StorageEnvironmentProd
	default:
		log.Fatalf("bad storage environment")
	}
	log.Println(config.mongoCollection, storageEnv)

	// create elwin config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ec, err := choices.NewChoices(
		ctx,
		choices.WithStorageConfig(config.mongoAddr, storageEnv),
		choices.UpdateInterval(time.Minute),
	)
	if err != nil {
		log.Fatal(err)
	}
	config.ec = ec
	config.readiness.storage = true

	http.Handle("/metrics", prometheus.Handler())

	go func() {
		config.readiness.httpServer = true
		config.ec.ErrChan <- http.ListenAndServe(config.jsonAddr, nil)
	}()

	go func() {
		lis, err := net.Listen("tcp", config.grpcAddr)
		if err != nil {
			config.ec.ErrChan <- fmt.Errorf("main: failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		elwin.RegisterElwinServer(grpcServer, &elwinServer{})
		config.readiness.grpcServer = true
		config.ec.ErrChan <- grpcServer.Serve(lis)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	config.readiness.errorHandler = true
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
	log.Printf("GetNamespaces: %v", id)
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

func logCloseErr(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("could not close response body: %s", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer jsonRequests.Inc()
	defer jsonDurations.Observe(float64(time.Since(start)))
	if err := r.ParseForm(); err != nil {
		log.Printf("could not parse form: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-type", "text/plain")
		if _, err := w.Write([]byte("invalid form data")); err != nil {
			log.Printf("could not write to root connection: %s", err)
		}
		return
	}
	if r.Body != nil {
		defer logCloseErr(r.Body)
	}

	var label string
	switch {
	case r.Form.Get("label") != "":
		label = r.Form.Get("label")
	case r.Form.Get("teamid") != "":
		label = r.Form.Get("teamid")
	case r.Form.Get("group-id") != "":
		label = r.Form.Get("group-id")
	default:
		label = ""
	}

	resp, err := config.ec.Namespaces(label, r.Form.Get("userid"))
	if err != nil {
		config.ec.ErrChan <- fmt.Errorf("rootHandler: couldn't get Namespaces: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := elwinJSON(resp, w); err != nil {
		config.ec.ErrChan <- err
		return
	}
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
	jsonComma    = []byte{',', ' '}
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
	ew.write([]byte{'\n'})
	if ew.err != nil {
		return ew.err
	}
	return nil
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("could not write to healthz connection: %s", err)
	}
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if config.readiness.storage && config.readiness.httpServer && config.readiness.grpcServer && config.readiness.errorHandler {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
