package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/foolusion/choices/storage/mongo"

	"gopkg.in/mgo.v2"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var err error
	if buf, err = httputil.DumpRequest(r, true); err != nil {
		log.Printf("could not dump request: %v", err)
		return
	}
	log.Printf("%s", buf)

	var result []mongo.MongoNamespace
	cfg.mongo.DB(cfg.mongoDB).C(cfg.testCollection).Find(nil).All(&result)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, result)
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
