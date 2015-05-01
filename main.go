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

var timeoutInMilliseconds = time.Millisecond * 100

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
		log.Println("get ["+k+"] in", timeSpent(t0), "ms")
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
	ttl, _ := strconv.Atoi(r.FormValue("ttl"))

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, k+": can't get value", http.StatusBadRequest)
		return
	}

	ch := make(chan error)
	go func() {
		ch <- ds.set(k, v, ttl)
	}()

	select {
	case err := <-ch:
		log.Println("set ["+k+"] in", timeSpent(t0), "ms")
		if err == nil {
			w.WriteHeader(http.StatusCreated)
		} else {
			switch err {
			case OVERSIZED_BODY, EMPTY_BODY:
				http.Error(w, k+": value is invalid", http.StatusBadRequest)
			case INVALID_KEY:
				http.Error(w, k+": invalid key", http.StatusBadRequest)
			default:
				http.Error(w, k+": cache server error", http.StatusInternalServerError)
			}
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
		ch <- ds.delete(k)
	}()

	select {
	case err := <-ch:
		log.Println("delete ["+k+"] in", timeSpent(t0), "ms")
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			switch err {
			case INVALID_KEY:
				http.Error(w, k+": invalid key", http.StatusBadRequest)
			default:
				http.Error(w, k+": cache server error", http.StatusInternalServerError)
			}
		}
	case <-time.After(timeoutInMilliseconds):
		returnTimeout(w, k)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	v, _ := ioutil.ReadAll(r.Body)
	ch := make(chan error)
	go func() {
		ch <- ds.append(k, v)
	}()

	select {
	case err := <-ch:
		log.Println("append ["+k+"] in", timeSpent(t0), "ms")
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			switch err {
			case NOT_FOUND_ERROR:
				http.Error(w, k+": not found", http.StatusNotFound)
			case OVERSIZED_BODY, EMPTY_BODY:
				http.Error(w, k+": value is invalid", http.StatusBadRequest)
			case INVALID_KEY:
				http.Error(w, k+": invalid key", http.StatusBadRequest)
			default:
				http.Error(w, k+": cache server error", http.StatusInternalServerError)
			}
		}
	case <-time.After(timeoutInMilliseconds):
		returnTimeout(w, k)
	}
}

func returnTimeout(w http.ResponseWriter, k string) {
	log.Println(k + ": timeout")
	http.Error(w, k+": timeout", http.StatusRequestTimeout)
}

func timeSpent(t0 int64) int64 {
	return int64(math.Floor(float64(time.Now().UnixNano()-t0)/1000000 + .5))
}
