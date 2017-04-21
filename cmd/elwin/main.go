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

	"github.com/Nordstrom/choices"
	"github.com/foolusion/elwinprotos/elwin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

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
	cfgStorageAddr = "storage_address"
	cfgJSONAddr    = "json_address"
	cfgGRPCAddr    = "grpc_address"
	cfgUInterval   = "update_interval"
	cfgRTimeout    = "read_timeout"
	cfgWTimeout    = "write_timeout"
	cfgITimeout    = "idle_timeout"
	cfgUFTimeout   = "update_fail_timeout"
	cfgProf        = "profiler"
)

func bind(s []string) error {
	if len(s) == 0 {
		return nil
	}
	if err := viper.BindEnv(s[0]); err != nil {
		return err
	}
	return bind(s[1:])
}

func main() {
	log.Println("Starting elwin...")

	viper.SetDefault(cfgStorageAddr, "elwin-storage:80")
	viper.SetDefault(cfgJSONAddr, ":8080")
	viper.SetDefault(cfgGRPCAddr, ":8081")
	viper.SetDefault(cfgUInterval, "10s")
	viper.SetDefault(cfgRTimeout, "5s")
	viper.SetDefault(cfgWTimeout, "5s")
	viper.SetDefault(cfgITimeout, "30s")
	viper.SetDefault(cfgUFTimeout, "15m")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/elwin")
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("no config file found")
		default:
			log.Fatalf("could not read config: %v", err)
		}
	}

	viper.SetEnvPrefix("elwin")
	if err := bind([]string{
		cfgStorageAddr,
		cfgJSONAddr,
		cfgGRPCAddr,
		cfgUInterval,
		cfgRTimeout,
		cfgWTimeout,
		cfgITimeout,
		cfgUFTimeout,
		"profiler",
	}); err != nil {
		log.Fatal(err)
	}

	// create elwin config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var interval, failTimeout time.Duration
	if ui, err := time.ParseDuration(viper.GetString(cfgUInterval)); err != nil {
		log.Fatal(err)
	} else if muft, err := time.ParseDuration(viper.GetString(cfgUFTimeout)); err != nil {
		log.Fatal(err)
	} else {
		interval, failTimeout = ui, muft
	}

	ec, err := choices.NewChoices(
		ctx,
		choices.WithGlobalSalt("choices"),
		choices.WithStorageConfig(viper.GetString(cfgStorageAddr), interval),
		choices.WithUpdateInterval(interval),
		choices.WithMaxUpdateFailTime(failTimeout),
	)
	if err != nil {
		log.Fatal(err)
	}

	// register prometheus metrics
	prometheus.MustRegister(jsonRequests)
	prometheus.MustRegister(jsonDurations)
	prometheus.MustRegister(updateErrors)

	ljson, err := net.Listen("tcp", viper.GetString(cfgJSONAddr))
	if err != nil {
		log.Fatalf("could not listen on %s: %v", viper.GetString(cfgJSONAddr), err)
	}
	defer ljson.Close()
	log.Printf("Listening for json on %s", viper.GetString(cfgJSONAddr))

	mux := http.NewServeMux()
	mux.Handle("/", &jsonServer{ec})
	mux.HandleFunc("/healthz", healthzHandler(map[string]interface{}{"storage": ec}))
	mux.HandleFunc("/readiness", healthzHandler(map[string]interface{}{"storage": ec}))
	mux.Handle("/metrics", promhttp.Handler())
	if viper.IsSet(cfgProf) {
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	}
	srv := http.Server{
		Handler: mux,
	}
	if rt, err := time.ParseDuration(viper.GetString(cfgRTimeout)); err != nil {
		log.Fatal(err)
	} else if wt, err := time.ParseDuration(viper.GetString(cfgWTimeout)); err != nil {
		log.Fatal(err)
	} else if it, err := time.ParseDuration(viper.GetString(cfgITimeout)); err != nil {
		log.Fatal(err)
	} else {
		srv.ReadTimeout = rt
		srv.WriteTimeout = wt
		srv.IdleTimeout = it
	}

	go func() {
		ec.ErrChan <- srv.Serve(ljson)
	}()

	go func() {
		lgrpc, err := net.Listen("tcp", viper.GetString(cfgGRPCAddr))
		if err != nil {
			ec.ErrChan <- fmt.Errorf("main: failed to listen: %v", err)
			return
		}
		defer lgrpc.Close()
		log.Printf("Listening for grpc on %s", viper.GetString(cfgGRPCAddr))

		grpcServer := grpc.NewServer()
		elwin.RegisterElwinServer(grpcServer, &elwinServer{ec})
		ec.ErrChan <- grpcServer.Serve(lgrpc)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-ec.ErrChan:
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

type elwinServer struct {
	*choices.Config
}

func (e *elwinServer) Get(ctx context.Context, req *elwin.GetRequest) (*elwin.GetReply, error) {
	if req == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Get: request is nil")
	}
	// TODO: we really need to pass in the requirements in the request.
	// This requires an update to elwin.proto
	selector, err := labels.Parse(req.Query)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse query")
	}
	resp, err := e.Experiments(req.UserID, selector)
	if err != nil {
		return nil, fmt.Errorf("error evaluating experiments for %s, %s: %v", req.Query, req.UserID, err)
	}

	if req.By != "" {
		byResp := make(map[string]*elwin.ExperimentList, 10)
		for _, v := range resp {
			if group, ok := v.Labels[req.By]; !ok {
				appendToGroup(byResp, v, "None")
			} else {
				appendToGroup(byResp, v, group)
			}
		}
		return &elwin.GetReply{
			Group: byResp,
		}, nil
	}

	exp := &elwin.GetReply{
		Experiments: make([]*elwin.Experiment, len(resp)),
	}

	for i, v := range resp {
		exp.Experiments[i] = &elwin.Experiment{
			Name:      v.Name,
			Namespace: v.Namespace,
			Labels:    v.Labels,
			Params:    make([]*elwin.Param, len(v.Params)),
		}

		for j, p := range v.Params {
			exp.Experiments[i].Params[j] = &elwin.Param{
				Name:  p.Name,
				Value: p.Value,
			}
		}
	}
	return exp, nil
}

func appendToGroup(br map[string]*elwin.ExperimentList, e choices.ExperimentResponse, group string) {
	if br[group] == nil {
		br[group] = &elwin.ExperimentList{}
	}
	elist := br[group].Experiments
	ee := &elwin.Experiment{
		Name:      e.Name,
		Namespace: e.Namespace,
		Labels:    e.Labels,
		Params:    make([]*elwin.Param, len(e.Params)),
	}
	for i, p := range e.Params {
		ee.Params[i] = &elwin.Param{
			Name:  p.Name,
			Value: p.Value,
		}
	}
	elist = append(elist, ee)
	br[group].Experiments = elist
}

func logCloseErr(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("could not close response body: %s", err)
	}
}

type jsonServer struct {
	*choices.Config
}

func (j *jsonServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
				j.ErrChan <- errors.Wrap(err, "could not create selection requirement")
				return
			}
			sel = sel.Add(*r)
		default:
			r, err := labels.NewRequirement(k, selection.In, v)
			if err != nil {
				j.ErrChan <- errors.Wrap(err, "could not create selection requirement")
				return
			}
			sel = sel.Add(*r)
		}
	}

	resp, err := j.Experiments(r.Form.Get("userid"), sel)
	if err != nil {
		j.ErrChan <- fmt.Errorf("rootHandler: couldn't get Experiments: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := elwinJSON(resp, w); err != nil {
		j.ErrChan <- err
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
