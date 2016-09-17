package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foolusion/choices"
)

func getUserIDs(variants map[string]string, value choices.Value, buf []byte, n, e, p string) map[string]string {
	found := 0
	for i := 1; i < len(buf); i++ {
		userID := fmt.Sprintf("%x", buf[:i])
		v, err := genValues(value, n, e, p, userID)
		if err != nil {
			return variants
		}
		key := fmt.Sprintf("%s.%s.%s.%s", n, e, p, v)
		if _, ok := variants[key]; !ok {
			found++
			variants[key] = userID
		}
		switch v := value.(type) {
		case *choices.Weighted:
			if found == len(v.Choices) {
				return variants
			}
		case *choices.Uniform:
			if found == len(v.Choices) {
				return variants
			}
		}
	}
	return variants
}

func genValues(v choices.Value, namespace, experiment, param, userID string) (string, error) {
	h, err := config.ec.HashExperience(namespace, experiment, param, userID)
	if err != nil {
		return "", err
	}
	return v.Value(h)
}

func genHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	r.ParseForm()
	namespaces := config.ec.Storage.Read()
	var teamNs []choices.Namespace
	for _, namespace := range namespaces {
		for _, label := range namespace.TeamID {
			if label == r.Form.Get("label") {
				teamNs = append(teamNs, namespace)
			}
		}
	}

	buf := make([]byte, 256)
	if _, err := rand.Read(buf); err != nil {
		return
	}

	variantURL := make(map[string]string, 100)
	for _, namespace := range teamNs {
		for _, experiment := range namespace.Experiments {
			for _, param := range experiment.Params {
				variantURL = getUserIDs(variantURL, param.Value, buf, namespace.Name, experiment.Name, param.Name)
			}
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(variantURL); err != nil {
		return
	}
}
