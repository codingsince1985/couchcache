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
	var err error
	url, bucket, pass := parseFlag()
	if ds, err = newDatastore(url, bucket, pass); err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	kr := r.PathPrefix("/key/{key}").Subrouter()
	kr.Methods("GET").HandlerFunc(getHandler)
	kr.Methods("POST").HandlerFunc(postHandler)
	kr.Methods("DELETE").HandlerFunc(deleteHandler)
	kr.Methods("PUT").HandlerFunc(putHandler)

	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func parseFlag() (string, *string, *string) {
	host := flag.String("host", "localhost", "host name (defaults to localhost)")
	port := flag.Int("port", 8091, "port number (defaults to 8091)")
	bucket := flag.String("bucket", "couchcache", "bucket name (defaults to couchcache)")
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
		returnTimeout(w, k)
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
		returnTimeout(w, k)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]
	ch := make(chan error)
	go func() {
		err := ds.delete(k)
		ch <- err
	}()

	select {
	case err := <-ch:
		if err == nil {
			log.Println("delete ["+k+"] in ", time.Now().UnixNano()-t0)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Println(err)
			http.Error(w, k+": cache server error", http.StatusInternalServerError)
		}
	case <-time.After(timeoutInMilliseconds):
		returnTimeout(w, k)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(k + ": bad request")
		http.Error(w, k+": bad request", http.StatusBadRequest)
		return
	}

	ch := make(chan error)
	go func() {
		err := ds.append(k, v)
		ch <- err
	}()

	select {
	case err = <-ch:
		if err == nil {
			log.Println("append ["+k+"] in ", time.Now().UnixNano()-t0)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Println(err)
			http.Error(w, k+": cache server error", http.StatusInternalServerError)
		}
	case <-time.After(timeoutInMilliseconds):
		returnTimeout(w, k)
	}
}

func returnTimeout(w http.ResponseWriter, k string) {
	log.Println(k + ": timeout")
	http.Error(w, k+": timeout", http.StatusRequestTimeout)
}
