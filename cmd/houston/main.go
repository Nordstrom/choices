package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/foolusion/choices"
	"github.com/foolusion/choices/storage/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	rootEndpoint      = "/"
	healthEndpoint    = "/healtz"
	readinessEndpoint = "readiness"
	launchPrefix      = "/launch/"
	deletePrefix      = "/delete/"
)

func init() {
	http.HandleFunc(rootEndpoint, rootHandler)
	http.HandleFunc(launchPrefix, launchHandler)
	http.HandleFunc(deletePrefix, deleteHandler)
	http.HandleFunc(healthEndpoint, healthHandler)
	http.HandleFunc(readinessEndpoint, readinessHandler)
}

type config struct {
	mongoAddr      string
	mongoDB        string
	testCollection string
	prodCollection string
	username       string
	password       string
	addr           string
	mongo          *mgo.Session
}

var cfg = config{
	mongoAddr:      "elwin-storage",
	mongoDB:        "elwin",
	testCollection: "test",
	prodCollection: "prod",
	username:       "elwin",
	password:       "philologist",
	addr:           ":8080",
}

const (
	envMongoAddress        = "MONGO_ADDRESS"
	envMongoDatabase       = "MONGO_DATABASE"
	envMongoTestCollection = "MONGO_TEST_COLLECTION"
	envMongoProdCollection = "MONGO_PROD_COLLECTION"
	envUsername            = "USERNAME"
	envPassword            = "PASSWORD"
	envAddr                = "ADDRESS"
)

func main() {
	log.Println("Starting Houston...")

	if os.Getenv(envMongoAddress) != "" {
		cfg.mongoAddr = os.Getenv(envMongoAddress)
		log.Printf("Setting Mongo Address: %q", cfg.mongoAddr)
	}
	if os.Getenv(envMongoDatabase) != "" {
		cfg.mongoDB = os.Getenv(envMongoDatabase)
		log.Printf("Setting Mongo Database: %q", cfg.mongoDB)
	}
	if os.Getenv(envMongoTestCollection) != "" {
		cfg.testCollection = os.Getenv(envMongoTestCollection)
		log.Printf("Setting Mongo Test Collection: %q", cfg.testCollection)
	}
	if os.Getenv(envMongoProdCollection) != "" {
		cfg.prodCollection = os.Getenv(envMongoProdCollection)
		log.Printf("Setting Mongo Prod Collection: %q", cfg.prodCollection)
	}
	if os.Getenv(envUsername) != "" {
		cfg.username = os.Getenv(envUsername)
		log.Printf("Setting Username: %q", cfg.username)
	}
	if os.Getenv(envPassword) != "" {
		cfg.password = os.Getenv(envPassword)
		log.Printf("Setting Password: %q", cfg.password)
	}

	errCh := make(chan error, 1)

	// setup mongo
	go func(c *config) {
		var err error
		c.mongo, err = mgo.Dial(c.mongoAddr)
		if err != nil {
			log.Printf("could not dial mongo database: %s", err)
			errCh <- err
		}
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

// Namespace container for data from mongo.
type Namespace struct {
	Name        string
	Labels      []string `bson:"teamid"`
	Experiments []struct {
		Name   string
		Params []mongo.Param
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
	TestRaw []Namespace
	ProdRaw []Namespace
	Test    []TableData
	Prod    []TableData
}

func namespaceToTableData(ns []Namespace) []TableData {
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
				switch p.Type {
				case choices.ValueTypeUniform:
					var uniform choices.Uniform
					p.Value.Unmarshal(&uniform)
					params[k].Values = strings.Join(uniform.Choices, ", ")
				case choices.ValueTypeWeighted:
					var weighted choices.Weighted
					p.Value.Unmarshal(&weighted)
					params[k].Values = strings.Join(weighted.Choices, ", ")
				}
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

	var test []Namespace
	cfg.mongo.DB(cfg.mongoDB).C(cfg.testCollection).Find(nil).All(&test)
	var prod []Namespace
	cfg.mongo.DB(cfg.mongoDB).C(cfg.prodCollection).Find(nil).All(&prod)

	data := rootTmplData{
		TestRaw: test,
		ProdRaw: prod,
		Test:    namespaceToTableData(test),
		Prod:    namespaceToTableData(prod),
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
<table>
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
	<th><a href="/delete/{{$exp.Name}}">Delete</a></th>
	<th><a href="/launch/{{$exp.Name}}">Launch</a></th>
</tr>
{{end}}
{{end}}
</table>
{{end}}

{{with .Prod}}
<h2>Prod</h2>
<table>
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
	<th><a href="/delete/{{$exp.Name}}">Delete</a></th>
	<th><a href="/launch/{{$exp.Name}}">Launch</a></th>
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
	experiment := r.URL.Path[len(launchPrefix):]

	// get the namespace from test
	test, err := mongo.QueryOne(cfg.mongo.DB(cfg.mongoDB).C(cfg.testCollection), bson.M{"experiments.name": experiment})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("not found"))
		return
	}
	var exp choices.Experiment
	for _, v := range test.Experiments {
		if v.Name == experiment {
			exp = v
			break
		}
	}

	// check for namespace in prod
	prod, err := mongo.QueryOne(cfg.mongo.DB(cfg.mongoDB).C(cfg.prodCollection), bson.M{"name": test.Name})
	if err == mgo.ErrNotFound {
		newProd := choices.Namespace{Name: test.Name, TeamID: test.TeamID, Experiments: []choices.Experiment{exp}}
		copy(newProd.Segments[:], choices.SegmentsAll[:])
		if err := newProd.Segments.Remove(&exp.Segments); err != nil {
			// this should never happen
			log.Println(err)
		}
		if err := mongo.Upsert(cfg.mongo.DB(cfg.mongoDB).C(cfg.prodCollection), newProd.Name, newProd); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("error launching to prod"))
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("something went wrong"))
		return
	}

	// subtract segments from prod namespace and add experiment
	if err := prod.Segments.Remove(&exp.Segments); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("not found"))
		return
	}
	prod.Experiments = append(prod.Experiments, exp)
	if err := mongo.Upsert(cfg.mongo.DB(cfg.mongoDB).C(cfg.prodCollection), prod.Name, prod); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("error launching to prod"))
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if err := cfg.mongo.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Not Ready"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}
