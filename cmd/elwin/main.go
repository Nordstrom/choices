package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/foolusion/choices"
)

func init() {
	http.HandleFunc("/", rootHandler)
}

func main() {
	log.Println("Starting elwin...")
	ns, err := choices.NewNamespace("t1", "test", []string{"userid"})
	if err != nil {
		log.Fatal("%v", err)
	}
	ns.Addexp(
		"aTest",
		[]choices.Param{{Name: "a", Value: &choices.Uniform{Choices: []string{"b", "c"}}}},
		128,
	)
	if err := choices.Addns(ns); err != nil {
		log.Fatalf("%v", err)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := choices.Namespaces(context.Background(), nil, "test", map[string][]string{"userid": []string{"some-user-id"}})
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	fmt.Fprintf(w, "%v", *resp)
}
