package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var ds datastore
var err error

func main() {
	ds, err = newDatastore("http://Administrator:password@localhost:8091/", "default")
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	kr := r.PathPrefix("/{key}").Subrouter()
	kr.Methods("GET").HandlerFunc(getHandler)
	kr.Methods("POST").HandlerFunc(postHandler)

	http.ListenAndServe(":8080", r)
}

func getHandler(rw http.ResponseWriter, r *http.Request) {
	maps := mux.Vars(r)
	key := maps["key"]
	if value := ds.get(key); value != nil {
		rw.Write(value)
	}
}

func postHandler(rw http.ResponseWriter, r *http.Request) {
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	maps := mux.Vars(r)
	key := maps["key"]
	ds.set(key, value)
}
