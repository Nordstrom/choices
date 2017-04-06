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
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/foolusion/elwinprotos/storage"
	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	cfgDBFile      = "db_file"
	cfgListenAddr  = "listen_address"
	cfgMetricsAddr = "metrics_address"
	cfgUser        = "user"
	cfgPassword    = "password"
	cfgSQLConnStr  = "sql_conn_str"
	cfgMailAddr    = "mail_address"
	cfgMailFrom    = "mail_from"
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
	log.Println("Starting bolt-store...")

	viper.SetDefault(cfgDBFile, "test.db")
	viper.SetDefault(cfgListenAddr, ":8080")
	viper.SetDefault(cfgMetricsAddr, ":8081")
	viper.SetDefault(cfgUser, "elwin")
	viper.SetDefault(cfgPassword, "")
	viper.SetDefault(cfgSQLConnStr, "localhost/elwin")
	viper.SetDefault(cfgMailAddr, "localhost:22")
	viper.SetDefault(cfgMailFrom, "elwin@nordstrom.com")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/elwin/bolt-store")
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("no config file found")
		default:
			log.Fatalf("could not read config: %v", err)
		}
	}

	viper.SetEnvPrefix("bolt_store")
	if err := bind([]string{
		cfgDBFile,
		cfgListenAddr,
		cfgMetricsAddr,
		cfgUser,
		cfgPassword,
		cfgSQLConnStr,
		cfgMailAddr,
		cfgMailFrom,
	}); err != nil {
		log.Fatal(err)
	}

	server, err := newServer(viper.GetString(cfgDBFile), viper.GetString(cfgUser), viper.GetString(cfgPassword), viper.GetString(cfgSQLConnStr))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("lisening for grpc on %q", viper.GetString(cfgListenAddr))
	lis, err := net.Listen("tcp", viper.GetString(cfgListenAddr))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	storage.RegisterElwinStorageServer(s, server)
	grpc_prometheus.Register(s)
	go func() {
		http.Handle("/metrics", prometheus.Handler())
		log.Printf("listening for /metrics on %q", viper.GetString(cfgMetricsAddr))
		log.Fatal(http.ListenAndServe(viper.GetString(cfgMetricsAddr), nil))
	}()

	log.Fatal(s.Serve(lis))
}

var (
	environmentStaging    = []byte("staging")
	environmentProduction = []byte("production")
)

type server struct {
	db       *bolt.DB
	redshift *sql.DB
}

func newServer(file, rUser, rPassword, rConnStr string) (*server, error) {
	switch {
	case rUser == "":
		return nil, errors.New("need to supply a user for sql connection")
	case rPassword == "":
		return nil, errors.New("need to supply a password for sql connection")
	case rConnStr == "":
		return nil, errors.New("need to supply a connection string for sql")
	}
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket(environmentStaging); err != nil {
			if err != bolt.ErrBucketExists {
				return err
			}
		}
		if _, err := tx.CreateBucket(environmentProduction); err != nil {
			if err != bolt.ErrBucketExists {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	connstr := fmt.Sprintf("postgres://%s:%s@%s?sslmode=require", rUser, rPassword, rConnStr)
	redshift, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	return &server{db: db, redshift: redshift}, nil
}

func (s *server) Close() error {
	s.redshift.Close()
	return s.db.Close()
}

// All returns all the namespaces for a given environment.
func (s *server) All(ctx context.Context, r *storage.AllRequest) (*storage.AllReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	ar := &storage.AllReply{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(env).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ns storage.Namespace
			if err := proto.Unmarshal(v, &ns); err != nil {
				return err
			}
			ar.Namespaces = append(ar.Namespaces, &ns)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ar, nil
}

// Create creates a namespace in the given environment.
func (s *server) Create(ctx context.Context, r *storage.CreateRequest) (*storage.CreateReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	ns := r.Namespace
	if ns == nil {
		return nil, fmt.Errorf("namespace is nil")
	}

	pns, err := proto.Marshal(ns)
	if err != nil {
		return nil, err
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(env).Put([]byte(ns.Name), pns)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.CreateReply{Namespace: ns}, nil
}

// Read returns the namespace matching the supplied name from the given
// environment.
func (s *server) Read(ctx context.Context, r *storage.ReadRequest) (*storage.ReadReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	if len(r.Name) == 0 {
		return nil, fmt.Errorf("name is empty")
	}

	ns := storage.Namespace{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket(env).Get([]byte(r.Name))
		if buf == nil {
			return grpc.Errorf(codes.NotFound, "key not found")
		}
		if err := proto.Unmarshal(buf, &ns); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.ReadReply{Namespace: &ns}, nil
}

// Update replaces the namespace in the given environment with the namespace
// supplied.
func (s *server) Update(ctx context.Context, r *storage.UpdateRequest) (*storage.UpdateReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	ns := r.GetNamespace()
	if ns == nil {
		return nil, fmt.Errorf("namespace is nil")
	}

	pns, err := proto.Marshal(ns)
	if err != nil {
		return nil, err
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(env).Put([]byte(ns.Name), pns)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &storage.UpdateReply{Namespace: ns}, nil
}

// Delete deletes the namespace from the given environment.
func (s *server) Delete(ctx context.Context, r *storage.DeleteRequest) (*storage.DeleteReply, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}
	env := envFromStorageRequest(r.Environment)

	if len(r.Name) == 0 {
		return nil, fmt.Errorf("name is empty")
	}

	ns := storage.Namespace{}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		buf := tx.Bucket(env).Get([]byte(r.Name))
		if buf == nil {
			return grpc.Errorf(codes.NotFound, "key not found")
		}
		if err := proto.Unmarshal(buf, &ns); err != nil {
			return err
		}
		return tx.Bucket(env).Delete([]byte(r.Name))
	}); err != nil {
		return nil, err
	}
	return &storage.DeleteReply{Namespace: &ns}, nil
}

const (
	insertElwinExperiment = `INSERT into elwin.experiment
(
	experiment_id,
	experiment_name,
	experiment_desc,
	platform,
	userid,
	programmanagerid,
	productmanagerid,
	team_id,
	hypothesis,
	kpi,
	timebound,
	plannedstarttime,
	plannedendtime,
	actionplannegative,
	actionplanneutral,
	experimenttype,
	namespace,
	segment,
	load_time
)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	insertElwinExperimentDetails = `INSERT into elwin.experiment_details
(
	experiment_id,
	experiment_name,
	paramname,
	choice,
	weight,
	label_key,
	label_value,
	load_time
)
VALUES
( $1, $2, $3, $4, $5, $6, $7, $8)`
)

func randName(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var str string
	b := make([]byte, 8)
	for i := 0; i < n; i++ {
		if _, err := rand.Read(b); err != nil {
			return "", errors.Wrap(err, "could not read from rand")
		}
		a := binary.BigEndian.Uint64(b) % uint64(len(alphabet))
		str += alphabet[a : a+1]
	}
	return str, nil
}

func uuid() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "could not read from rand")
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func (s *server) ExperimentIntake(ctx context.Context, r *storage.ExperimentIntakeRequest) (*storage.ExperimentIntakeReply, error) {

	var err error
	if r.Namespace.Experiments[0].Name, err = randName(7); err != nil {
		return nil, errors.Wrap(err, "failed creating experiment name")
	}

	tx, err := s.redshift.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "could not begin a transaction")
	}

	platforms := strings.Split(r.Namespace.Experiments[0].Labels["platform"], ",")
	for _, p := range platforms {
		if err := insertExperiment(ctx, tx, s.Create, p, r); err != nil {
			return nil, errors.Wrapf(err, "could not insert experiment for platform %s", p)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit the experiment intake transaction")
	}
	r.Namespace.Experiments[0].Labels["platform"] = strings.Join(platforms, ", ")
	if err := sendEmail(r); err != nil {
		log.Println(err)
	}

	return &storage.ExperimentIntakeReply{}, nil
}

func insertExperiment(
	ctx context.Context,
	tx *sql.Tx,
	create func(context.Context, *storage.CreateRequest) (*storage.CreateReply, error),
	platform string,
	r *storage.ExperimentIntakeRequest,
) error {
	r.Namespace.Experiments[0].Labels["platform"] = platform
	var err error
	if r.Namespace.Name, err = randName(7); err != nil {
		return errors.Wrap(err, "could not generate namespace name")
	} else if r.Namespace.Experiments[0].Id, err = uuid(); err != nil {
		return errors.Wrap(err, "could not generate experiment uuid")
	}

	if err := processOrRollback(insertElwinExperimentTx, tx, r); err != nil {
		return errors.Wrap(err, "insert to experiment failed")
	} else if err := processOrRollback(insertElwinExperimentDetailsTx, tx, r); err != nil {
		return errors.Wrap(err, "insert to experiment_details failed")
	}

	if _, err := create(ctx, &storage.CreateRequest{
		Environment: storage.Staging,
		Namespace:   r.Namespace,
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "failed to rollback transaction")
		}
		return errors.Wrap(err, "could not create namespace in bolt-store")
	}

	return nil
}

func processOrRollback(f func(tx *sql.Tx, r *storage.ExperimentIntakeRequest) error, tx *sql.Tx, r *storage.ExperimentIntakeRequest) error {
	if err := f(tx, r); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "could not rollback transaction")
		}
		return errors.Wrap(err, "could not process request")
	}
	return nil
}

func insertElwinExperimentTx(tx *sql.Tx, r *storage.ExperimentIntakeRequest) error {
	var buf bytes.Buffer
	wc := base64.NewEncoder(base64.StdEncoding, &buf)
	wc.Write(r.Namespace.Experiments[0].Segments)
	wc.Close()

	var endDate string
	if d := r.Metadata.PlannedEndTime; d == "" {
		endDate = "9999-01-01"
	} else {
		endDate = d
	}

	if res, err := tx.Exec(
		insertElwinExperiment,
		r.Namespace.Experiments[0].Id,
		r.Namespace.Experiments[0].Name,
		r.Namespace.Experiments[0].DetailName,
		r.Namespace.Experiments[0].Labels["platform"],
		r.Metadata.UserID,
		r.Metadata.ProgramManagerID,
		r.Metadata.ProductManagerID,
		r.Namespace.Experiments[0].Labels["team"],
		r.Metadata.Hypothesis,
		r.Metadata.Kpi,
		r.Metadata.TimeBound,
		r.Metadata.PlannedStartTime,
		endDate,
		r.Metadata.ActionPlanNegative,
		r.Metadata.ActionPlanNeutral,
		r.Metadata.ExperimentType,
		r.Namespace.Name,
		buf.String(),
		time.Now().Format("2006-01-02 15:04:05.000 MST"),
	); err != nil {
		return errors.Wrap(err, "failed exec insert into elwin.experiment")
	} else if i, err := res.RowsAffected(); err != nil {
		return errors.Wrap(err, "could not get RowsAffected")
	} else if i == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func insertElwinExperimentDetailsTx(tx *sql.Tx, r *storage.ExperimentIntakeRequest) error {
	for _, e := range r.Namespace.Experiments {
		for _, p := range e.Params {
			for ci, c := range p.Value.Choices {
				for lk, lv := range e.Labels {
					var weight float64
					if len(p.Value.Weights) == 0 {
						weight = 0
					} else {
						weight = p.Value.Weights[ci]
					}

					if res, err := tx.Exec(
						insertElwinExperimentDetails,
						e.Id,
						e.Name,
						p.Name,
						c,
						weight,
						lk,
						lv,
						time.Now().Format("2006-01-02 15:04:05.000 MST"),
					); err != nil {
						return errors.Wrap(err, "failed exec insert to elwin.experiment_details")
					} else if i, err := res.RowsAffected(); err != nil {
						return errors.Wrap(err, "could not get Rows Affected")
					} else if i == 0 {
						return errors.Wrap(err, "no rows affected")
					}
				}
			}
		}
	}
	return nil
}

type param struct {
	Name    string
	Choices []choice
}

type choice struct {
	Name   string
	Weight float64
}

func sendEmail(r *storage.ExperimentIntakeRequest) error {
	to := []string{
		r.Metadata.UserID + "@nordstrom.com",
		r.Metadata.ProductManagerID + "@nordstrom.com",
		r.Metadata.ProgramManagerID + "@nordstrom.com",
	}
	cc := []string{
		viper.GetString(cfgMailFrom),
	}

	var data = struct {
		ToAddr             string
		CcAddr             string
		ExperimentName     string
		Params             []param
		Team               string
		Platform           string
		UserID             string
		ProgramManagerID   string
		ProductManagerID   string
		Hypothesis         string
		ActionPlanNeutral  string
		ActionPlanNegative string
	}{
		ToAddr:             strings.Join(to, ", "),
		CcAddr:             strings.Join(cc, ", "),
		ExperimentName:     r.Namespace.Experiments[0].DetailName,
		Params:             makeParam(r.Namespace.Experiments[0].Params),
		Team:               r.Namespace.Experiments[0].Labels["team"],
		Platform:           r.Namespace.Experiments[0].Labels["platform"],
		UserID:             r.Metadata.UserID,
		ProgramManagerID:   r.Metadata.ProgramManagerID,
		ProductManagerID:   r.Metadata.ProductManagerID,
		Hypothesis:         r.Metadata.Hypothesis,
		ActionPlanNeutral:  r.Metadata.ActionPlanNeutral,
		ActionPlanNegative: r.Metadata.ActionPlanNegative,
	}

	crlf := []byte("\r\n")
	var msg bytes.Buffer
	fmt.Fprintf(&msg, "To: %s\r\n", data.ToAddr)
	fmt.Fprintf(&msg, "Cc: %s\r\n", data.CcAddr)
	fmt.Fprintf(&msg, "Subject: =?utf-8?Q?Experiment Intake for Experiment Name: %s?=\r\n", data.ExperimentName)
	msg.WriteString("MIME-Version: 1.0\r\n")
	mw := multipart.NewWriter(&msg)
	fmt.Fprintf(&msg, "Content-Type: multipart/mixed;\r\n\tboundary=%q\r\n\r\n", mw.Boundary())
	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=UTF-8")
	pw, err := mw.CreatePart(textHeader)
	if err != nil {
		return errors.Wrap(err, "could not create multipart header")
	}
	if err := mailTmpl.Execute(pw, data); err != nil {
		return errors.Wrap(err, "could not execute mail template")
	}
	pw.Write(crlf)
	pw.Write(crlf)
	if err := mw.Close(); err != nil {
		return errors.Wrap(err, "could not close ")
	}

	return sendMail(viper.GetString(cfgMailAddr), noAuth{}, viper.GetString(cfgMailFrom), append(to, cc...), msg.Bytes())
}

func makeParam(params []*storage.Param) []param {
	var ps []param
	for _, p := range params {
		var choices []choice
		for i, c := range p.Value.Choices {
			switch len(p.Value.Weights) {
			case 0:
				choices = append(choices, choice{Name: c})
			default:
				choices = append(choices, choice{Name: c, Weight: p.Value.Weights[i]})
			}
		}
		ps = append(ps, param{Name: p.Name, Choices: choices})
	}
	return ps
}

var mailTmpl = template.Must(template.New("mailTmpl").Parse(mailTmplText))

const mailTmplText = `Thank you for your submission.  Your request has been created, and the TTO team has been notified.  TTO will evaluate the request and be in touch with your team's program manager and product manager with next steps.

NEO DATA

EXPERIMENT NAME:
{{ .ExperimentName }}
{{ range $ip, $p := .Params -}}
PARAM {{ $ip }}:
{{ $p.Name }}
{{ range $ic, $c := $p.Choices -}}
CHOICE {{ $ic }}:
{{ $c.Name }}{{ if ne $c.Weight 0.0 }} {{ $c.Weight }}{{ end }}
{{ end }}{{ end }}

EXPERIMENT METADATA

TEAM:
{{ .Team }}
PLATFORM:
{{ .Platform }}
REQUESTER LAN ID:
{{ .UserID }}
PROGRAM MANAGER ID:
{{ .ProgramManagerID }}
PRODUCT MANAGER ID:
{{ .ProductManagerID }}

HYPOTHESIS:
{{ .Hypothesis }}

ACTION PLAN FOR NEUTRAL RESULTS:
{{ .ActionPlanNeutral }}

ACTION PLAN FOR NEGATIVE RESULTS:
{{ .ActionPlanNegative }}`

func envFromStorageRequest(e storage.Environment) []byte {
	switch e {
	case storage.Staging:
		return environmentStaging
	case storage.Production:
		return environmentProduction
	default:
		return environmentStaging
	}
}
