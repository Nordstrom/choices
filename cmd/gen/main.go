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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/foolusion/choices"
	storage "github.com/foolusion/choices/elwinstorage"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	storageAddr         = "elwin-storage:80"
	listenAddr          = ":8080"
	genEndpoint         = "/"
	bookmarkletEndpoint = "/bookmarklet"
	globalSalt          = "choices"
	grpcTimeout         = 500 * time.Millisecond

	envJavascriptFile = "JAVASCRIPT_FILE"
	envURL            = "URL"
)

var (
	config = struct {
		cc     *grpc.ClientConn
		client storage.ElwinStorageClient
	}{}

	ErrNotFound = errors.New("could not generate a matching cookie value")
)

var bookmarkletTmpl *template.Template

func init() {
	http.HandleFunc(genEndpoint, genHandler)
	http.HandleFunc(bookmarkletEndpoint, bookmarkletHandler)
	choices.SetGlobalSalt(globalSalt)

}

func main() {
	var err error
	config.cc, err = grpc.Dial(storageAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	config.client = storage.NewElwinStorageClient(config.cc)
	f, err := os.Open(os.Getenv(envJavascriptFile))
	if err != nil {
		log.Fatal(err)
	}
	out, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	str := strings.Replace(string(out), `"`, `%22`, -1)
	bookmarkletTmpl = template.Must(template.New("bookmarklet").Parse(bookmarkletHTML))
	_ = template.Must(bookmarkletTmpl.New("javascript").Parse(str))

	log.Println(http.ListenAndServe(listenAddr, nil))
}

func bookmarkletHandler(w http.ResponseWriter, _ *http.Request) {
	if err := bookmarkletTmpl.Execute(w, os.Getenv(envURL)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var bookmarkletHTML = `<!doctype html>
<html lang="en">
<head>
<title>bookmarklets - gen</title>
</head>
<body>
<h1>Gen Bookmarklet</h1>
<p>Drag this link to the bookmark bar</p>
<a href="javascript:{{template "javascript" .}}">Elwin</a>
</body>
</html>
`

func genHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ar := &storage.AllRequest{Environment: storage.Environment_Production}
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	resp, err := config.client.All(ctx, ar)
	if err != nil {
		var errCode int
		switch grpc.Code(err) {
		case codes.Canceled, codes.DeadlineExceeded:
			errCode = http.StatusRequestTimeout
		case codes.InvalidArgument:
			errCode = http.StatusBadRequest
		default:
			errCode = http.StatusInternalServerError
		}
		http.Error(w, err.Error(), errCode)
		return
	}

	var namespaces []choices.Namespace
	for _, namespace := range resp.GetNamespaces() {
		cns, err := choices.FromNamespace(namespace)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		namespaces = append(namespaces, cns)
	}
	ev, err := gen(namespaces)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	if err := enc.Encode(ev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type experimentValue struct {
	NamespaceName  string             `json:"namespaceName"`
	ExperimentName string             `json:"experimentName"`
	Labels         string             `json:"labels"`
	Params         map[string][]param `json:"params"`
}

type param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func gen(ns []choices.Namespace) ([]experimentValue, error) {
	buf := make([]byte, 256)
	if _, err := rand.Read(buf); err != nil {
		return nil, errors.Wrap(err, "could not read random bytes")
	}
	var expVal []experimentValue
	for _, n := range ns {
		for _, e := range n.Experiments {
			ev := experimentValue{
				NamespaceName:  n.Name,
				ExperimentName: e.Name,
				Labels:         strings.Join(n.TeamID, ", "),
				Params:         make(map[string][]param, 16),
			}

			cookies, err := cookie(buf, n.Name, e)
			if err != nil {
				return nil, errors.Wrap(err, "could not generate cookie values")
			}

			for key, cookie := range cookies {
				ev.Params[cookie] = unkey(key)
			}
			expVal = append(expVal, ev)
		}
	}
	return expVal, nil
}

type paramKey struct {
	param
	more interface{}
}

func key(params ...param) paramKey {
	var more interface{}
	if len(params) > 1 {
		more = key(params[1:]...)
	}
	return paramKey{param: params[0], more: more}
}

func unkey(p paramKey) []param {
	var ps []param
	for {
		ps = append(ps, p.param)
		if p.more == nil {
			return ps
		}
		p = p.more.(paramKey)
	}
	return ps
}

func cookie(buf []byte, namespace string, experiment choices.Experiment) (map[paramKey]string, error) {
	num := uniqueParams(experiment.Params)

	cookies := make(map[paramKey]string, 16)
	for i := 1; i < len(buf); i++ {
		if len(cookies) == num {
			return cookies, nil
		}
		userID := hex.EncodeToString(buf[:i])
		if !choices.InSegment(namespace, userID, experiment.Segments) {
			continue
		}
		var paramKeys []param
		for _, p := range experiment.Params {
			val, err := genValues(p.Value, namespace, experiment.Name, p.Name, userID)
			if err != nil {
				return nil, errors.Wrap(err, "could not generate value")
			}
			paramKeys = append(paramKeys, param{Name: p.Name, Value: val})
		}
		k := key(paramKeys...)
		if _, ok := cookies[k]; !ok {
			cookies[k] = userID
		}
	}
	return nil, ErrNotFound
}

func uniqueParams(params []choices.Param) int {
	res := 1
	for _, param := range params {
		switch v := param.Value.(type) {
		case *choices.Uniform:
			res *= len(v.Choices)
		case *choices.Weighted:
			res *= len(v.Choices)
		}
	}

	return res
}

func genValues(v choices.Value, namespace, experiment, param, userID string) (string, error) {
	h, err := choices.HashExperience(namespace, experiment, param, userID)
	if err != nil {
		return "", err
	}
	return v.Value(h)
}
