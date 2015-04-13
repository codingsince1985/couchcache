package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

var ds datastore

func main() {
	ds = newDatastore("http://Administrator:password@localhost:8091/", "default")

	r := mux.NewRouter()
	r.HandleFunc("/{key}", handler)

	http.ListenAndServe(":8080", r)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	maps := mux.Vars(r)
	key := maps["key"]
	v := ds.get(key)
	rw.Write(v)
}
