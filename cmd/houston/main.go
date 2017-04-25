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
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	rootEndpoint      = "/"
	launchEndpoint    = "/launch"
	deleteEndpoint    = "/delete"
	healthEndpoint    = "/healthz"
	readinessEndpoint = "/readiness"

	cfgStorageAddr = "clients"
	cfgListenAddr  = "listen_address"
)

var (
	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "houston",
		Name:      "requests",
		Help:      "The number of requests to houston",
	}, []string{"type"})
	errRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "nordstrom",
		Subsystem: "houston",
		Name:      "errors",
		Help:      "The number of requests to houston",
	}, []string{"type"})
)

func init() {
	prometheus.MustRegister(requests)
	prometheus.MustRegister(errRequests)
}

type config struct {
	clients map[string]storageConfig
}

type storageConfig struct {
	client storage.ElwinStorageClient
	conn   *grpc.ClientConn
}

func main() {
	log.Println("Starting Houston...")
	viper.SetDefault(cfgStorageAddr, map[string]string{"dev": "elwin-storage:80", "prod": "elwin-storage:80"})
	viper.SetDefault(cfgListenAddr, ":8080")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/houston")
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("no config file found")
		default:
			log.Fatalf("could not read config: %v", err)
		}
	}

	http.Handle("/metrics", promhttp.Handler())

	errCh := make(chan error, 1)

	cfg := config{
		clients: make(map[string]storageConfig, 10),
	}

	for name, addr := range viper.GetStringMapString(cfgStorageAddr) {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		cfg.clients[name] = storageConfig{client: storage.NewElwinStorageClient(conn), conn: conn}
	}

	// set up http server
	lis, err := net.Listen("tcp", viper.GetString(cfgListenAddr))
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	log.Printf("Listening on %s", viper.GetString(cfgListenAddr))
	mux := http.NewServeMux()
	mux.Handle(rootEndpoint, rootHandler(cfg))
	mux.Handle(launchEndpoint, launchHandler(cfg))
	mux.Handle(deleteEndpoint, deleteHandler(cfg))
	mux.HandleFunc(healthEndpoint, healthHandler)
	mux.HandleFunc(readinessEndpoint, readinessHandler)
	server := http.Server{
		Handler: mux,
	}
	go func() {
		errCh <- server.Serve(lis)
	}()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				log.Println(err)
				for _, client := range cfg.clients {
					client.conn.Close()
				}
				func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					if err := server.Shutdown(ctx); err != nil {
						log.Println(err)
					}
					return
				}()
			}
		}
	}
}

// tableData container for data to be output.
type tableData struct {
	ID        string
	Name      string
	Namespace string
	Labels    string
	Params    []tableDataParam
}

type tableDataParam struct {
	Name   string
	Values string
}

type rootTmplData struct {
	DataRaw map[string][]*storage.Experiment
	Data    map[string]struct {
		OtherEnvs []string
		Exps      []tableData
	}
}

func namespaceToTableData(exps []*storage.Experiment) []tableData {
	tds := make([]tableData, len(exps))
	for i, v := range exps {
		tds[i] = tableData{
			ID:        v.Id,
			Name:      v.Name,
			Namespace: v.Namespace,
			Labels:    labels.Set(v.Labels).String(),
			Params:    make([]tableDataParam, len(v.Params)),
		}

		for j, p := range v.Params {
			tds[i].Params[j] = tableDataParam{
				Name:   p.Name,
				Values: strings.Join(p.Value.Choices, ", "),
			}
		}
	}
	return tds
}

func otherEnvs(c config, env string) []string {
	keys := make([]string, 0, len(c.clients))
	for k := range c.clients {
		if k != env {
			keys = append(keys, k)
		}
	}
	return keys
}

type rootHandler config

func (root rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requests.With(prometheus.Labels{"type": "dashboard"}).Inc()
	data := rootTmplData{Data: make(map[string]struct {
		OtherEnvs []string
		Exps      []tableData
	}, len(root.clients))}
	for name, client := range root.clients {
		resp, err := client.client.List(context.Background(), &storage.ListRequest{})
		if err != nil {
			errRequests.With(prometheus.Labels{"type": "dashboard"}).Inc()
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		data.Data[name] = struct {
			OtherEnvs []string
			Exps      []tableData
		}{
			OtherEnvs: otherEnvs(config(root), name),
			Exps:      namespaceToTableData(resp.Experiments),
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=uttf-8")
	if err := rootTemplate.Execute(w, data); err != nil {
		errRequests.With(prometheus.Labels{"type": "dashboard"}).Inc()
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
{{ range $name, $exps := .Data }}
<h2>{{ $name }}</h2>
<table class="table table-striped">
<tr>
  <th>Namespace</th>
  <th>Labels</th>
  <th>Experiment</th>
  <th>Params</th>
  <th>Delete?</th>
  <th>Launch?</th>
</tr>
{{ range $i, $exp := $exps.Exps -}}
<tr>
	<th>{{ $exp.Namespace }}</th>
	<th>{{ $exp.Labels }}</th>
	<th>{{ $exp.Name }}</th>
	<th>{{ range $exp.Params }}<strong>{{ .Name }}</strong>: ({{ .Values }})<br/>{{ end }}</th>
	<th><a href="/delete?id={{ $exp.ID }}&environment={{ $name }}">Delete</a></th>
	<th>
	{{- range $j, $env := $exps.OtherEnvs -}}
		{{- if $j }}, {{ end -}}
		<a href="/launch?id={{ $exp.ID }}&from={{ $name }}&to={{ $env }}">{{ $env }}</a>
	{{- end -}}
	</th>
</tr>
{{- end }}
</table>
{{- end }}
</div>
</body>
</html>
`

type launchHandler config

func (l launchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requests.With(prometheus.Labels{"type": "launch"}).Inc()
	if err := r.ParseForm(); err != nil {
		errRequests.With(prometheus.Labels{"type": "launch"}).Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getResp, err := l.clients[r.Form.Get("from")].client.Get(context.Background(), &storage.GetRequest{Id: r.Form.Get("id")})
	if err != nil {
		errRequests.With(prometheus.Labels{"type": "launch"}).Inc()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = l.clients[r.Form.Get("to")].client.Set(context.Background(), &storage.SetRequest{Experiment: getResp.Experiment})
	if err != nil {
		errRequests.With(prometheus.Labels{"type": "launch"}).Inc()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(struct {
		ID    string `json:"id"`
		State string `json:"state"`
	}{
		ID:    r.Form.Get("id"),
		State: "START",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", "http://"+r.Form.Get("to")+"/api/v1/experiment-change-state", buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	redshiftResp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer redshiftResp.Body.Close()
	log.Println(redshiftResp.Body)

	http.Redirect(w, r, "/", http.StatusFound)
}

type deleteHandler config

func (d deleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requests.With(prometheus.Labels{"type": "delete"}).Inc()
	err := r.ParseForm()
	if err != nil {
		errRequests.With(prometheus.Labels{"type": "delete"}).Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sc, ok := d.clients[r.Form.Get("environment")]
	if !ok {
		errRequests.With(prometheus.Labels{"type": "delete"}).Inc()
		http.Error(w, "bad environment", http.StatusBadRequest)
		return
	}

	_, err = sc.client.Remove(context.TODO(), &storage.RemoveRequest{
		Id: r.Form.Get("id"),
	})
	if err != nil {
		errRequests.With(prometheus.Labels{"type": "delete"}).Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	requests.With(prometheus.Labels{"type": "health"}).Inc()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	requests.With(prometheus.Labels{"type": "readiness"}).Inc()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}
