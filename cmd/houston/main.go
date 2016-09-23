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
)

const (
	rootEndpoint      = "/"
	healthEndpoint    = "/healthz"
	readinessEndpoint = "/readiness"
	launchPrefix      = "/launch"
	deletePrefix      = "/delete"
)

func init() {
	http.HandleFunc(rootEndpoint, rootHandler)
	http.HandleFunc(launchPrefix, launchHandler)
	http.HandleFunc(deletePrefix, deleteHandler)
	http.HandleFunc(healthEndpoint, healthHandler)
	http.HandleFunc(readinessEndpoint, readinessHandler)
}

type config struct {
	storageAddr string
	mongoDB     string
	username    string
	password    string
	addr        string
	conn        *grpc.ClientConn
	esc         storage.ElwinStorageClient
}

var cfg = config{
	storageAddr: "elwin-storage:80",
	mongoDB:     "elwin",
	username:    "elwin",
	password:    "philologist",
	addr:        ":8080",
}

const (
	envStorageAddress = "STORAGE_ADDRESS"
	envMongoDatabase  = "MONGO_DATABASE"
	envUsername       = "USERNAME"
	envPassword       = "PASSWORD"
	envAddr           = "ADDRESS"
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

func main() {
	log.Println("Starting Houston...")

	configFromEnv(&cfg.storageAddr, envStorageAddress)
	configFromEnv(&cfg.mongoDB, envMongoDatabase)
	configFromEnv(&cfg.username, envUsername)
	configFromEnv(&cfg.password, envPassword)
	configFromEnv(&cfg.addr, envAddr)

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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var err error
	if buf, err = httputil.DumpRequest(r, true); err != nil {
		log.Printf("could not dump request: %v", err)
		return
	}
	log.Printf("%s", buf)

	stagingReply, err := cfg.esc.All(context.TODO(), &storage.AllRequest{Environment: storage.Environment_Staging})
	if err != nil {
		log.Printf("AllRequest failed: %s", err)
	}
	productionReply, err := cfg.esc.All(context.TODO(), &storage.AllRequest{Environment: storage.Environment_Production})
	if err != nil {
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
	<th><a href="/delete?experiment={{$exp.Name}}">Delete</a></th>
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
	<th><a href="/delete?experiment={{$exp.Name}}">Delete</a></th>
	<th><a href="/launch?namespace={{$ns.Name}}&experiment={{$exp.Name}}">Launch</a></th>
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
	if err := r.ParseForm(); err != nil {
		logAndWriteError(err, "could not parse form", w, http.StatusBadRequest)
		return
	}
	namespace := r.Form.Get("namespace")
	experiment := r.Form.Get("experiment")

	// get the namespace from test
	stagingReply, err := cfg.esc.Read(
		context.TODO(),
		&storage.ReadRequest{Name: namespace, Environment: storage.Environment_Staging},
	)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("not found"))
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

	// check for namespace in prod
	productionReply, err := cfg.esc.Read(
		context.TODO(),
		&storage.ReadRequest{Name: namespace, Environment: storage.Environment_Production},
	)
	if err != nil {
		switch grpc.Code(err) {
		case codes.NotFound:
			createErr := createNamespace(ns.Name, ns.Labels, exp)
			if createErr != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Write([]byte("error launching to prod"))
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("something went wrong"))
			return
		}
	}

	prod := productionReply.GetNamespace()
	if prod == nil {
		return
	}

	prodNS, err := choices.FromNamespace(prod)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("something went wrong"))
		return
	}

	// subtract segments from prod namespace and add experiment
	if err := prodNS.Segments.Remove(&exp.Segments); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("not found"))
		return
	}
	prodNS.Experiments = append(prodNS.Experiments, exp)
	ureq := &storage.UpdateRequest{
		Namespace:   prodNS.ToNamespace(),
		Environment: storage.Environment_Production,
	}
	_, err = cfg.esc.Update(context.TODO(), ureq)
	if err != nil {
		logAndWriteError(err, "error launching to prod", w, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func createNamespace(name string, labels []string, exp choices.Experiment) error {
	newProd := choices.Namespace{Name: name, TeamID: labels, Experiments: []choices.Experiment{exp}}
	copy(newProd.Segments[:], choices.SegmentsAll[:])
	if err := newProd.Segments.Remove(&exp.Segments); err != nil {
		return errors.Wrap(err, "error removing segments, this should never happen...")
	}
	_, err := cfg.esc.Create(context.TODO(), &storage.CreateRequest{
		Namespace:   newProd.ToNamespace(),
		Environment: storage.Environment_Production,
	})

	return err
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func logAndWriteError(err error, errMsg string, w http.ResponseWriter, httpStatus int) {
	log.Println(errors.Wrap(err, errMsg))
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(errMsg))
}
