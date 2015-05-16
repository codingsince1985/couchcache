package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

var ds datastorer

var timeout = time.Millisecond * 100

func main() {
	if d, err := newDatastore(); err != nil {
		log.Fatalln(err)
	} else {
		ds = datastorer(d)
	}

	r := mux.NewRouter()
	kr := r.PathPrefix("/key/{key}").Subrouter()
	kr.Methods("GET").HandlerFunc(getHandler)
	kr.Methods("POST").HandlerFunc(postHandler)
	kr.Methods("DELETE").HandlerFunc(deleteHandler)
	kr.Methods("PUT").HandlerFunc(putHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	if err := ds.validKey(k); err != nil {
		http.Error(w, k+": invalid key", http.StatusBadRequest)
		return
	}

	ch := make(chan []byte)
	go func() {
		ch <- ds.get(k)
	}()

	select {
	case v := <-ch:
		if v != nil {
			log.Println("get ["+k+"] in", timeSpent(t0), "ms")
			w.Write(v)
		} else {
			log.Println(k + ": not found")
			http.Error(w, k+": not found", http.StatusNotFound)
		}
	case <-time.After(timeout):
		returnTimeout(w, k)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]
	ttl, _ := strconv.Atoi(r.FormValue("ttl"))

	if v, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, k+": can't get value", http.StatusBadRequest)
		return
	} else {
		if err = ds.validKey(k); err == nil {
			if err = ds.validValue(v); err == nil {
				go func() {
					ds.set(k, v, ttl)
				}()

				log.Println("set ["+k+"] in", timeSpent(t0), "ms")
				w.WriteHeader(http.StatusCreated)
				return
			}
		}
		datastoreErrorToHTTPError(err, w)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	if err := ds.delete(k); err == nil {
		log.Println("delete ["+k+"] in", timeSpent(t0), "ms")
		w.WriteHeader(http.StatusNoContent)
	} else {
		datastoreErrorToHTTPError(err, w)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	if v, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, k+": can't get value", http.StatusBadRequest)
		return
	} else {
		if err = ds.append(k, v); err == nil {
			log.Println("append ["+k+"] in", timeSpent(t0), "ms")
			w.WriteHeader(http.StatusOK)
		} else {
			datastoreErrorToHTTPError(err, w)
		}
	}
}

func returnTimeout(w http.ResponseWriter, k string) {
	log.Println(k + ": timeout")
	http.Error(w, k+": timeout", http.StatusRequestTimeout)
}

func timeSpent(t0 int64) int64 {
	return int64(math.Floor(float64(time.Now().UnixNano()-t0)/1000000 + .5))
}

func datastoreErrorToHTTPError(err error, w http.ResponseWriter) {
	switch err {
	case errNotFound:
		http.Error(w, "key not found", http.StatusNotFound)
	case errEmptyBody:
		http.Error(w, "empty value", http.StatusBadRequest)
	case errOversizedBody:
		http.Error(w, "oversized value", http.StatusBadRequest)
	case errInvalidKey:
		http.Error(w, "invalid key", http.StatusBadRequest)
	case errKeyExists:
		http.Error(w, "key exists", http.StatusBadRequest)
	default:
		http.Error(w, "cache server error", http.StatusInternalServerError)
	}
}
