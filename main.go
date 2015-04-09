package main

import (
	"github.com/couchbase/go-couchbase"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	bucket := "default"
	c, _ := couchbase.Connect("http://Administrator:password@localhost:8091/")
	p, _ := c.GetPool("default")
	b, _ := p.GetBucketWithAuth(bucket, bucket, "")
	v, _ := b.GetRaw("key")
	rw.Write(v)
}
