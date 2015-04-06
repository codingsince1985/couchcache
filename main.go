package main

import "net/http"

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("hello"))
}
