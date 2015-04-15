package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var ds datastore

func main() {
	if d, err := newDatastore("http://Administrator:password@localhost:8091/", "default"); err != nil {
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

func getHandler(w http.ResponseWriter, r *http.Request) {
	k := mux.Vars(r)["key"]
	if v := ds.get(k); v != nil {
		w.Write(v)
	} else {
		http.Error(w, k+" not found", http.StatusNotFound)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	k := mux.Vars(r)["key"]
	ds.set(k, v)
	w.WriteHeader(http.StatusCreated)
}
