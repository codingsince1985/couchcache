package main

import (
	"github.com/couchbase/go-couchbase"
	"github.com/gorilla/mux"
	"net/http"
)

var c couchbase.Client
var p couchbase.Pool
var b *couchbase.Bucket

const bucket = "default"

func main() {
	c, _ = couchbase.Connect("http://Administrator:password@localhost:8091/")
	p, _ = c.GetPool("default")
	b, _ = p.GetBucketWithAuth(bucket, bucket, "")

	r := mux.NewRouter()
	r.HandleFunc("/id/{key}", handler)

	http.ListenAndServe(":8080", r)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	v, _ := b.GetRaw(key)
	rw.Write(v)
}
