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
	"context"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/foolusion/choices"
	storage "github.com/foolusion/choices/elwinstorage"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type config struct {
	storageAddr string
	mongoDB     string
	username    string
	password    string
	addr        string
	conn        *grpc.ClientConn
	esc         storage.ElwinStorageClient
}

const (
	rootEndpoint      = "/"
	launchEndpoint    = "/launch"
	deleteEndpoint    = "/delete"
	healthEndpoint    = "/healthz"
	readinessEndpoint = "/readiness"

	envStorageAddress = "STORAGE_ADDRESS"
	envMongoDatabase  = "MONGO_DATABASE"
	envUsername       = "USERNAME"
	envPassword       = "PASSWORD"
	envAddr           = "ADDRESS"

	mbr  = "bad request"
	mint = "internal"
	mnf  = "not found"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")

	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "houston",
		Name:      "requests",
		Help:      "The number of requests to houston",
	}, []string{"path"})
	errRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "houston",
		Name:      "error_requests",
		Help:      "The number of requests to houston",
	}, []string{"path", "type"})

	cfg = config{
		storageAddr: "elwin-storage:80",
		mongoDB:     "elwin",
		username:    "elwin",
		password:    "philologist",
		addr:        ":8080",
	}
)

func configFromEnv(dst *string, env string) {
	if dst == nil {
		return
	}
	if os.Getenv(env) != "" {
		*dst = os.Getenv(env)
		log.Printf("Setting %s: %q", env, *dst)
	}
}

func init() {
	http.HandleFunc(rootEndpoint, rootHandler)
	http.HandleFunc(launchEndpoint, launchHandler)
	http.HandleFunc(deleteEndpoint, deleteHandler)
	http.HandleFunc(healthEndpoint, healthHandler)
	http.HandleFunc(readinessEndpoint, readinessHandler)
	prometheus.MustRegister(requests)
	prometheus.MustRegister(errRequests)
}

func main() {
	log.Println("Starting Houston...")

	configFromEnv(&cfg.storageAddr, envStorageAddress)
	configFromEnv(&cfg.mongoDB, envMongoDatabase)
	configFromEnv(&cfg.username, envUsername)
	configFromEnv(&cfg.password, envPassword)
	configFromEnv(&cfg.addr, envAddr)

	http.Handle("/metrics", prometheus.Handler())

	errCh := make(chan error, 1)

	// setup grpc
	go func(c *config) {
		var err error
		c.conn, err = grpc.Dial(cfg.storageAddr, grpc.WithInsecure())
		if err != nil {
			log.Printf("could not dial grpc storage: %s", err)
			errCh <- err
		}
		c.esc = storage.NewElwinStorageClient(c.conn)
	}(&cfg)

	go func() {
		errCh <- http.ListenAndServe(cfg.addr, nil)
	}()
	for {
		select {
		case err := <-errCh:
			if err != nil {
				log.Fatal(err)
				// graceful shutdown
				return
			}
		}
	}
}

// TableData container for data to be output.
type TableData struct {
	Name        string
	Labels      string
	Experiments []struct {
		Name   string
		Params []struct {
			Name   string
			Values string
		}
	}
}

type rootTmplData struct {
	TestRaw []*storage.Namespace
	ProdRaw []*storage.Namespace
	Test    []TableData
	Prod    []TableData
}

func namespaceToTableData(ns []*storage.Namespace) []TableData {
	tableData := make([]TableData, len(ns))
	for i, v := range ns {
		tableData[i].Name = v.Name
		tableData[i].Labels = strings.Join(v.Labels, ", ")
		experiments := make(
			[]struct {
				Name   string
				Params []struct {
					Name   string
					Values string
				}
			}, len(v.Experiments))
		tableData[i].Experiments = experiments
		for j, e := range v.Experiments {
			tableData[i].Experiments[j].Name = e.Name
			params := make(
				[]struct {
					Name   string
					Values string
				}, len(e.Params))
			for k, p := range e.Params {
				params[k].Name = p.Name
				params[k].Values = strings.Join(p.Value.Choices, ", ")
			}
			tableData[i].Experiments[j].Params = params
		}
	}
	return tableData
}

func incErrMetrics(err error, labels prometheus.Labels) {
	switch grpc.Code(err) {
	case codes.NotFound:
		labels["type"] = "not found"
		errRequests.With(labels).Inc()
	case codes.InvalidArgument:
		labels["type"] = "bad request"
		errRequests.With(labels).Inc()
	case codes.Internal:
		labels["type"] = "internal"
		errRequests.With(labels).Inc()
	default:
		labels["type"] = "internal"
		errRequests.With(labels).Inc()
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	labelGen := l(rootEndpoint)
	requests.With(prometheus.Labels{"path": rootEndpoint}).Inc()
	var buf []byte
	var err error
	if buf, err = httputil.DumpRequest(r, true); err != nil {
		errRequests.With(labelGen(mint)).Inc()
		log.Printf("could not dump request: %v", err)
		return
	}
	log.Printf("%s", buf)

	stagingReply, err := cfg.esc.All(context.TODO(), &storage.AllRequest{Environment: storage.Environment_Staging})
	if err != nil {
		incErrMetrics(err, labelGen(mint))
		log.Printf("AllRequest failed: %s", err)
	}
	productionReply, err := cfg.esc.All(context.TODO(), &storage.AllRequest{Environment: storage.Environment_Production})
	if err != nil {
		incErrMetrics(err, labelGen(mint))
		log.Printf("AllRequest failed: %s", err)
	}

	data := rootTmplData{
		TestRaw: stagingReply.GetNamespaces(),
		ProdRaw: productionReply.GetNamespaces(),
		Test:    namespaceToTableData(stagingReply.GetNamespaces()),
		Prod:    namespaceToTableData(productionReply.GetNamespaces()),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := rootTemplate.Execute(w, data); err != nil {
		errRequests.With(labelGen(mint)).Inc()
		log.Println(err)
	}
}

var rootTemplate = template.Must(template.New("root").Parse(rootTmpl))

const rootTmpl = `<!doctype html>
<html lang="en">
<head>
<title>Houston!</title>
</head>
<body>
<h1>Houston</h1>
<div>
{{with .Test}}
<h2>Test</h2>
<table class="table table-striped">
<tr>
  <th>Namespace</th>
  <th>Labels</th>
  <th>Experiment</th>
  <th>Params</th>
  <th>Delete?</th>
  <th>Launch?</th>
</tr>
{{range $ns := .}}
{{range $exp := .Experiments}}
<tr>
	<th>{{$ns.Name}}</th>
	<th>{{$ns.Labels}}</th>
	<th>{{$exp.Name}}</th>
	<th>{{range .Params}}<strong>{{.Name}}</strong>: ({{.Values}})<br/>{{end}}</th>
	<th><a href="/delete?namespace={{$ns.Name}}&experiment={{$exp.Name}}&environment=staging">Delete</a></th>
	<th><a href="/launch?namespace={{$ns.Name}}&experiment={{$exp.Name}}">Launch</a></th>
</tr>
{{end}}
{{end}}
</table>
{{end}}

{{with .Prod}}
<h2>Prod</h2>
<table class="table table-striped">
<tr>
  <th>Namespace</th>
  <th>Labels</th>
  <th>Experiment</th>
  <th>Params</th>
  <th>Delete?</th>
  <th>Launch?</th>
</tr>
{{range $ns := .}}
{{range $exp := .Experiments}}
<tr>
	<th>{{$ns.Name}}</th>
	<th>{{$ns.Labels}}</th>
	<th>{{$exp.Name}}</th>
	<th>{{range .Params}}<strong>{{.Name}}</strong>: ({{.Values}})<br/>{{end}}</th>
	<th><a href="/delete?namespace={{$ns.Name}}&experiment={{$exp.Name}}&environment=production">Delete</a></th>
</tr>
{{end}}
{{end}}
</table>
{{end}}

</div>
</body>
</html>
`

func launchHandler(w http.ResponseWriter, r *http.Request) {
	labelGen := l(deleteEndpoint)
	requests.With(prometheus.Labels{"path": launchEndpoint}).Inc()
	log.Println("starting launch...")
	if err := r.ParseForm(); err != nil {
		errRequests.With(labelGen(mbr)).Inc()
		logAndWriteError(err, "could not parse form", w, http.StatusBadRequest)
		return
	}
	namespace := r.Form.Get("namespace")
	experiment := r.Form.Get("experiment")

	log.Println("reading staging namespace...")
	// get the namespace from test
	stagingReply, err := cfg.esc.Read(
		context.TODO(),
		&storage.ReadRequest{Name: namespace, Environment: storage.Environment_Staging},
	)
	if err != nil {
		incErrMetrics(err, labelGen(mint))
		logAndWriteError(err, "not found", w, http.StatusNotFound)
		return
	}
	var exp choices.Experiment
	ns := stagingReply.GetNamespace()
	if ns == nil {
		return
	}
	for _, v := range ns.Experiments {
		if v.Name == experiment {
			exp = choices.FromExperiment(v)
			break
		}
	}

	log.Println("reading production namespace...")
	// check for namespace in prod
	productionReply, err := cfg.esc.Read(
		context.TODO(),
		&storage.ReadRequest{Name: namespace, Environment: storage.Environment_Production},
	)
	if err != nil {
		switch grpc.Code(err) {
		case codes.NotFound:
			log.Println("not found in production")
			createErr := createNamespace(ns.Name, ns.Labels, exp)
			if createErr != nil {
				errRequests.With(labelGen(mint)).Inc()
				logAndWriteError(err, "error launching to prod", w, http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		default:
			errRequests.With(labelGen(mint)).Inc()
			logAndWriteError(err, "something went wrong", w, http.StatusInternalServerError)
			return
		}
	}

	prod := productionReply.GetNamespace()
	if prod == nil {
		errRequests.With(labelGen(mint)).Inc()
		logAndWriteError(err, "something went wrong", w, http.StatusInternalServerError)
		return
	}

	log.Println(prod, exp)

	prodNS, err := choices.FromNamespace(prod)
	if err != nil {
		errRequests.With(labelGen(mint)).Inc()
		logAndWriteError(err, "something went wrong", w, http.StatusInternalServerError)
		return
	}

	// subtract segments from prod namespace and add experiment
	seg, err := prodNS.Segments.Claim(exp.Segments)
	if err != nil {
		errRequests.With(labelGen(mnf)).Inc()
		logAndWriteError(err, "not found", w, http.StatusNotFound)
		return
	}
	prodNS.Segments = seg

	prodNS.Experiments = append(prodNS.Experiments, exp)
	log.Println(prod)
	ureq := &storage.UpdateRequest{
		Namespace:   prodNS.ToNamespace(),
		Environment: storage.Environment_Production,
	}
	log.Println(ureq)
	_, err = cfg.esc.Update(context.TODO(), ureq)
	if err != nil {
		errRequests.With(labelGen(mint)).Inc()
		logAndWriteError(err, "error launching to prod", w, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func createNamespace(name string, labels []string, exp choices.Experiment) error {
	log.Println("starting create namespace")
	newProd := choices.Namespace{Name: name, Labels: labels, Experiments: []choices.Experiment{exp}}
	cr, err := cfg.esc.Create(context.TODO(), &storage.CreateRequest{
		Namespace:   newProd.ToNamespace(),
		Environment: storage.Environment_Production,
	})
	if err != nil {
		return err
	}

	log.Println(*cr, err)

	return nil
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	labelGen := l(deleteEndpoint)
	err := r.ParseForm()
	if err != nil {
		errRequests.With(labelGen(mbr)).Inc()
		logAndWriteError(err, "could not parse form", w, http.StatusBadRequest)
	}
	requests.With(prometheus.Labels{"path": deleteEndpoint}).Inc()

	var storageEnv storage.Environment
	switch r.Form.Get("environment") {
	case "staging":
		storageEnv = storage.Environment_Staging
	case "production":
		storageEnv = storage.Environment_Production
	default:
		storageEnv = storage.Environment_Staging
	}

	prodReadReq, err := cfg.esc.Read(context.TODO(), &storage.ReadRequest{
		Name:        r.Form.Get("namespace"),
		Environment: storage.Environment_Production,
	})

	var prodNS choices.Namespace
	if err != nil {
		prodNS = choices.Namespace{}
	} else {
		prodNS, err = choices.FromNamespace(prodReadReq.Namespace)
		if err != nil {
			errRequests.With(labelGen(mint)).Inc()
			logAndWriteError(err, "could not parse namespace", w, http.StatusInternalServerError)
			return
		}
	}

	prodIndex := -1
	for i, exp := range prodNS.Experiments {
		if exp.Name == r.Form.Get("experiment") {
			prodIndex = i
			break
		}
	}

	if prodIndex >= 0 && storageEnv == storage.Environment_Production {
		if err := deleteExperiment(prodNS, storageEnv, prodIndex); err != nil {
			errRequests.With(labelGen(mint)).Inc()
			logAndWriteError(err, "could not delete prod experiment", w, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if prodIndex >= 0 && storageEnv == storage.Environment_Staging {
		errRequests.With(labelGen(mbr)).Inc()
		logAndWriteError(ErrBadRequest, "test still in prod", w, http.StatusBadRequest)
		return
	} else if prodIndex < 0 && storageEnv == storage.Environment_Production {
		errRequests.With(labelGen(mnf)).Inc()
		logAndWriteError(ErrNotFound, "test is not in prod", w, http.StatusNotFound)
		return
	}

	var stagingNS choices.Namespace
	stagReadReq, err := cfg.esc.Read(context.TODO(), &storage.ReadRequest{
		Name:        r.Form.Get("namespace"),
		Environment: storage.Environment_Staging,
	})
	if err != nil {
		http.Error(w, "namespace not found: "+err.Error(), http.StatusNotFound)
		return
	} else {
		stagingNS, err = choices.FromNamespace(stagReadReq.Namespace)
		if err != nil {
			errRequests.With(labelGen(mint)).Inc()
			logAndWriteError(err, "could not parse staging namespace", w, http.StatusInternalServerError)
			return
		}
	}

	stagIndex := -1
	for i, exp := range stagingNS.Experiments {
		if exp.Name == r.Form.Get("experiment") {
			stagIndex = i
			break
		}
	}
	if stagIndex == -1 {
		http.Error(w, "could not match experiment", http.StatusNotFound)
		return
	}

	if err := deleteExperiment(stagingNS, storageEnv, stagIndex); err != nil {
		errRequests.With(labelGen(mint)).Inc()
		logAndWriteError(err, "could not delete experiment", w, http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func deleteExperiment(ns choices.Namespace, env storage.Environment, index int) error {
	if len(ns.Experiments) == 1 {
		if _, err := cfg.esc.Delete(context.TODO(), &storage.DeleteRequest{
			Name:        ns.Name,
			Environment: env,
		}); err != nil {
			return err
		}
		return nil
	}
	ns.Experiments[index] = ns.Experiments[len(ns.Experiments)-1]
	ns.Experiments[len(ns.Experiments)-1] = choices.Experiment{}
	ns.Experiments = ns.Experiments[:len(ns.Experiments)-1]
	if _, err := cfg.esc.Update(context.TODO(), &storage.UpdateRequest{
		Namespace:   ns.ToNamespace(),
		Environment: env,
	}); err != nil {
		return err
	}
	return nil
}

func logAndWriteError(err error, errMsg string, w http.ResponseWriter, httpStatus int) {
	log.Println(errors.Wrap(err, errMsg))
	http.Error(w, errMsg, httpStatus)
}

func l(p string) func(string) prometheus.Labels {
	return func(t string) prometheus.Labels {
		return prometheus.Labels{"path": p, "type": t}
	}
}
