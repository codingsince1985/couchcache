package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var ds datastore

var timeoutInMilliseconds = time.Millisecond * 100

func main() {
	url, bucket, pass := parseFlag()
	if d, err := newDatastore(url, bucket, pass); err != nil {
		log.Fatalln(err)
	} else {
		ds = d
	}

	r := mux.NewRouter()
	kr := r.PathPrefix("/key/{key}").Subrouter()
	kr.Methods("GET").HandlerFunc(getHandler)
	kr.Methods("POST").HandlerFunc(postHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func parseFlag() (string, *string, *string) {
	host := flag.String("host", "localhost", "host name (defaults to localhost)")
	port := flag.Int("port", 8091, "port (defaults to 8091)")
	bucket := flag.String("bucket", "couchcache", "username (defaults to couchcache)")
	pass := flag.String("pass", "password", "password (defaults to password)")

	flag.Parse()

	url := fmt.Sprintf("http://%s:%d", *host, *port)
	log.Println(url)
	return url, bucket, pass
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]
	ch := make(chan []byte)
	go func() {
		v := ds.get(k)
		ch <- v
	}()

	select {
	case v := <-ch:
		log.Println("get ["+k+"] in ", time.Now().UnixNano()-t0)
		if v != nil {
			w.Write(v)
		} else {
			log.Println(k + ": not found")
			http.Error(w, k+": not found", http.StatusNotFound)
		}
	case <-time.After(timeoutInMilliseconds):
		log.Println(k + ": timeout")
		http.Error(w, k+": timeout", http.StatusRequestTimeout)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]
	ttl, err := strconv.Atoi(r.FormValue("ttl"))
	if err != nil {
		ttl = 0
	}

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(k + ": bad request")
		http.Error(w, k+": bad request", http.StatusBadRequest)
		return
	}

	ch := make(chan error)
	go func() {
		err := ds.set(k, v, ttl)
		ch <- err
	}()

	select {
	case err = <-ch:
		if err == nil {
			log.Println("set ["+k+"] in ", time.Now().UnixNano()-t0)
			w.WriteHeader(http.StatusCreated)
		} else {
			log.Println(err)
			http.Error(w, k+": cache server error", http.StatusInternalServerError)
		}
	case <-time.After(timeoutInMilliseconds):
		log.Println(k + ": timeout")
		http.Error(w, k+": timeout", http.StatusRequestTimeout)
	}
}
