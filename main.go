package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var ds datastore

func main() {
	url, bucket := parseFlag()
	if d, err := newDatastore(url, *bucket); err != nil {
		log.Fatalln(err)
	} else {
		ds = d
	}

	r := mux.NewRouter()
	kr := r.PathPrefix("/key/{key}").Subrouter()
	kr.Methods("GET").HandlerFunc(getHandler)
	kr.Methods("POST").HandlerFunc(postHandler)

	http.ListenAndServe(":8080", r)
}

func parseFlag() (string, *string) {
	user := flag.String("user", "Administrator", "username (defaults to Administrator)")
	pass := flag.String("pass", "password", "password (defaults to password)")
	host := flag.String("host", "localhost", "host name (defaults to localhost)")
	port := flag.Int("port", 8091, "port (defaults to 8091)")
	bucket := flag.String("bucket", "default", "port (defaults to default)")

	flag.Parse()

	url := fmt.Sprintf("http://%s:%s@%s:%d", *user, *pass, *host, *port)
	log.Println(url)
	return url, bucket
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	k := mux.Vars(r)["key"]
	if v := ds.get(k); v != nil {
		w.Write(v)
	} else {
		http.Error(w, k+" not found", http.StatusNotFound)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	ttl, err := strconv.Atoi(r.FormValue("ttl"))
	if err != nil {
		ttl = 0
	}

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	k := mux.Vars(r)["key"]
	ds.set(k, v, ttl)
	w.WriteHeader(http.StatusCreated)
}
