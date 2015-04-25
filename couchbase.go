package main

import (
	"flag"
	"fmt"
	"github.com/couchbase/go-couchbase"
	"log"
)

const MAX_TTL_IN_SEC = 60 * 60 * 24 * 30

type couchbaseDatastore couchbase.Bucket

func newDatastore() (ds datastorer, err error) {
	url, bucket, pass := parseFlag()

	c, err := couchbase.ConnectWithAuthCreds(url, bucket, pass)
	if err != nil {
		return ds, err
	}

	p, err := c.GetPool("default")
	if err != nil {
		return ds, err
	}

	b, err := p.GetBucketWithAuth(bucket, bucket, pass)
	if err != nil {
		return ds, err
	}

	return (datastorer)((*couchbaseDatastore)(b)), nil
}

func parseFlag() (string, string, string) {
	host := flag.String("host", "localhost", "host name (defaults to localhost)")
	port := flag.Int("port", 8091, "port number (defaults to 8091)")
	bucket := flag.String("bucket", "couchcache", "bucket name (defaults to couchcache)")
	pass := flag.String("pass", "password", "password (defaults to password)")

	flag.Parse()

	url := fmt.Sprintf("http://%s:%d", *host, *port)
	log.Println(url)
	return url, *bucket, *pass
}

func (ds *couchbaseDatastore) get(k string) []byte {
	v, err := (*couchbase.Bucket)(ds).GetRaw(k)
	if err != nil {
		return nil
	}

	return []byte(v)
}

func (ds *couchbaseDatastore) set(k string, v []byte, ttl int) error {
	if ttl > MAX_TTL_IN_SEC {
		ttl = MAX_TTL_IN_SEC
	}

	err := (*couchbase.Bucket)(ds).SetRaw(k, ttl, v)
	return err
}

func (ds *couchbaseDatastore) delete(k string) error {
	err := (*couchbase.Bucket)(ds).Delete(k)
	return err
}

func (ds *couchbaseDatastore) append(k string, v []byte) error {
	err := (*couchbase.Bucket)(ds).Append(k, v)
	return err
}
