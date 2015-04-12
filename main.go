package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

var p *pool

func main() {
	p = newPool("http://Administrator:password@localhost:8091/")

	r := mux.NewRouter()
	r.HandleFunc("/{bucket}/{key}", handler)

	http.ListenAndServe(":8080", r)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	maps := mux.Vars(r)
	bucket := maps["bucket"]
	b, err := p.getBucket(bucket)
	if err != nil {
		rw.Write([]byte(""))
	} else {
		key := maps["key"]
		v, _ := b.GetRaw(key)
		rw.Write(v)
	}
}
