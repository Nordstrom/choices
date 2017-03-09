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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"golang.org/x/net/context"

	"github.com/Nordstrom/choices"
	"github.com/foolusion/elwinprotos/elwin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var config = struct {
	ec              *choices.Config
	grpcAddr        string
	jsonAddr        string
	mongoAddr       string
	mongoDB         string
	mongoCollection string
	readTimeout     string
	writeTimeout    string
	idleTimeout     string
	readiness       struct {
		storage      bool
		grpcServer   bool
		errorHandler bool
		httpServer   bool
	}
	updateInterval    string
	maxUpdateFailTime string
}{
	jsonAddr:          ":8081",
	grpcAddr:          ":8080",
	mongoAddr:         "elwin-storage:80",
	mongoDB:           "elwin",
	mongoCollection:   "test",
	readTimeout:       "5s",
	writeTimeout:      "5s",
	idleTimeout:       "30s",
	updateInterval:    "10s",
	maxUpdateFailTime: "5m",
}

var (
	jsonRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "elwin",
		Name:      "json_requests",
		Help:      "The number of json requests received.",
	})
	jsonDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "nordstrom",
		Subsystem: "elwin",
		Name:      "json_durations_nanoseconds",
		Help:      "json latency distributions.",
		Buckets:   prometheus.ExponentialBuckets(1, 10, 10),
	})
	updateErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "elwin",
		Name:      "update_errors",
		Help:      "The number of errors while updating storage.",
	})
)

const (
	envJSONAddr          = "JSON_ADDRESS"
	envGRPCAddr          = "GRPC_ADDRESS"
	envMongoAddr         = "MONGO_ADDRESS"
	envMongoDB           = "MONGO_DATABASE"
	envMongoColl         = "MONGO_COLLECTION"
	envReadTimeout       = "READ_TIMEOUT"
	envWriteTimeout      = "WRITE_TIMEOUT"
	envIdleTimeout       = "IDLE_TIMEOUT"
	envProfiler          = "PROFILER"
	envUpdateInterval    = "UPDATE_INTERVAL"
	envMaxUpdateFailTime = "MAX_UPDATE_FAIL_TIME"
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
	env(&config.readTimeout, envReadTimeout)
	env(&config.writeTimeout, envWriteTimeout)
	env(&config.idleTimeout, envIdleTimeout)
	env(&config.updateInterval, envUpdateInterval)
	env(&config.maxUpdateFailTime, envMaxUpdateFailTime)

	var storageEnv int
	switch config.mongoCollection {
	case "staging", "dev", "test":
		storageEnv = choices.StorageEnvironmentDev
	case "production", "prod":
		storageEnv = choices.StorageEnvironmentProd
	default:
		log.Fatal("bad storage environment")
	}
	log.Println(config.mongoCollection, storageEnv)

	// create elwin config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var updateInterval, maxUpdateFailTime time.Duration
	if ui, err := time.ParseDuration(config.updateInterval); err != nil {
		log.Fatal(err)
	} else if muft, err := time.ParseDuration(config.maxUpdateFailTime); err != nil {
		log.Fatal(err)
	} else {
		updateInterval, maxUpdateFailTime = ui, muft
	}

	ec, err := choices.NewChoices(
		ctx,
		choices.WithGlobalSalt("choices"),
		choices.WithStorageConfig(config.mongoAddr, storageEnv, updateInterval),
		choices.WithUpdateInterval(updateInterval),
		choices.WithMaxUpdateFailTime(maxUpdateFailTime),
	)
	if err != nil {
		log.Fatal(err)
	}
	config.ec = ec
	config.readiness.storage = true

	// register prometheus metrics
	prometheus.MustRegister(jsonRequests)
	prometheus.MustRegister(jsonDurations)
	prometheus.MustRegister(updateErrors)

	ljson, err := net.Listen("tcp", config.jsonAddr)
	if err != nil {
		log.Fatalf("could not listen on %s: %v", config.jsonAddr, err)
	}
	defer ljson.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthzHandler(map[string]interface{}{"storage": config.ec}))
	mux.HandleFunc("/readiness", readinessHandler)
	mux.Handle("/metrics", prometheus.Handler())
	if len(os.Getenv(envProfiler)) > 0 {
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	}
	srv := http.Server{
		Handler: mux,
	}
	if rt, err := time.ParseDuration(config.readTimeout); err != nil {
		log.Fatal(err)
	} else if wt, err := time.ParseDuration(config.writeTimeout); err != nil {
		log.Fatal(err)
	} else if it, err := time.ParseDuration(config.idleTimeout); err != nil {
		log.Fatal(err)
	} else {
		srv.ReadTimeout = rt
		srv.WriteTimeout = wt
		srv.IdleTimeout = it
	}

	go func() {
		config.readiness.httpServer = true
		config.ec.ErrChan <- srv.Serve(ljson)
	}()

	go func() {
		lgrpc, err := net.Listen("tcp", config.grpcAddr)
		if err != nil {
			config.ec.ErrChan <- fmt.Errorf("main: failed to listen: %v", err)
			return
		}
		defer lgrpc.Close()

		grpcServer := grpc.NewServer()
		elwin.RegisterElwinServer(grpcServer, &elwinServer{})
		config.readiness.grpcServer = true
		config.ec.ErrChan <- grpcServer.Serve(lgrpc)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	config.readiness.errorHandler = true
	for {
		select {
		case err := <-config.ec.ErrChan:
			switch errors.Cause(err).(type) {
			case choices.ErrUpdateStorage:
				updateErrors.Inc()
			}
			log.Println(err)
		case s := <-signalChan:
			log.Printf("Captured %v. Exitting...", s)
			cancel()
			ctx, sdcancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer sdcancel()
			srv.Shutdown(ctx)
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
		return nil, grpc.Errorf(codes.InvalidArgument, "GetNamespaces: no Identifier received")
	}
	// TODO: we really need to pass in the requirements in the request. This requires an update to elwin.proto
	selector := labels.NewSelector()
	r, err := labels.NewRequirement("team", selection.In, []string{id.TeamID})
	if err != nil {
		return nil, errors.Wrap(err, "could not create requirement")
	}
	selector = selector.Add(*r)
	resp, err := config.ec.Namespaces(id.UserID, selector)
	if err != nil {
		return nil, fmt.Errorf("error resolving namespaces for %s, %s: %v", id.TeamID, id.UserID, err)
	}

	exp := &elwin.Experiments{
		Experiments: make(map[string]*elwin.Experiment, len(resp)),
	}

	for _, v := range resp {
		exp.Experiments[v.Name] = &elwin.Experiment{
			Namespace: v.Namespace,
			Params:    make([]*elwin.Param, len(v.Params)),
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
	defer func() {
		jsonDurations.Observe(float64(time.Since(start)))
	}()
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

	sel := labels.NewSelector()

	for k, v := range r.Form {
		switch k {
		case "userid":
			continue
		case "team", "label", "teamid", "group-id":
			r, err := labels.NewRequirement("team", selection.In, v)
			if err != nil {
				config.ec.ErrChan <- errors.Wrap(err, "could not create selection requirement")
				return
			}
			sel = sel.Add(*r)
		default:
			r, err := labels.NewRequirement(k, selection.In, v)
			if err != nil {
				config.ec.ErrChan <- errors.Wrap(err, "could not create selection requirement")
				return
			}
			sel = sel.Add(*r)
		}
	}

	resp, err := config.ec.Namespaces(r.Form.Get("userid"), sel)
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

func (ew *errWriter) writeKey(buf []byte) {
	ew.write(jsonQuote)
	ew.write(buf)
	ew.write(jsonQuote)
	ew.write(jsonKeyEnd)
}

func (ew *errWriter) writeString(buf []byte) {
	ew.write(jsonQuote)
	ew.write(buf)
	ew.write(jsonQuote)
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
	ew.writeKey([]byte("experiments"))
	ew.write(jsonOpenObj)
	for i, exp := range er {
		ew.writeKey([]byte(exp.Name))
		ew.write(jsonOpenObj)
		ew.writeKey([]byte("namespace"))
		ew.writeString([]byte(exp.Namespace))
		ew.write(jsonComma)
		ew.writeKey([]byte("params"))
		ew.write(jsonOpenObj)
		for j, param := range exp.Params {
			ew.writeKey([]byte(param.Name))
			ew.writeString([]byte(param.Value))
			if j < len(exp.Params)-1 {
				ew.write(jsonComma)
			}
		}
		ew.write(jsonCloseObj)
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

func healthzHandler(healthChecks map[string]interface{}) http.HandlerFunc {
	type healthy interface {
		IsHealthy() error
	}

	return func(w http.ResponseWriter, r *http.Request) {
		errs := make(map[string]string, len(healthChecks))
		for key, healthChecker := range healthChecks {
			if hc, ok := healthChecker.(healthy); ok {
				err := hc.IsHealthy()
				if err != nil {
					errs[key] = err.Error()
				}
			}
		}
		if len(errs) != 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			if err := enc.Encode(errs); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("could not write to healthz connection: %s", err)
		}
	}
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	if config.readiness.storage && config.readiness.httpServer && config.readiness.grpcServer && config.readiness.errorHandler {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
