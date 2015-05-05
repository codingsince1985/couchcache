package main

import (
	"flag"
	"fmt"
	"github.com/couchbase/go-couchbase"
	"github.com/couchbase/gomemcached"
	"log"
)

const (
	MAX_TTL_IN_SEC   = 60 * 60 * 24 * 30
	MAX_SIZE_IN_BYTE = 20 * 1024 * 1024
	MAX_KEY_LENGTH   = 250
)

type couchbaseDatastore couchbase.Bucket

func newDatastore() (ds *couchbaseDatastore, err error) {
	url, bucket, pass := parseFlag()

	if c, err := couchbase.ConnectWithAuthCreds(url, bucket, pass); err == nil {
		if p, err := c.GetPool("default"); err == nil {
			if b, err := p.GetBucketWithAuth(bucket, bucket, pass); err == nil {
				return (*couchbaseDatastore)(b), nil
			}
		}
	}
	return nil, err
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
	if v, err := (*couchbase.Bucket)(ds).GetRaw(k); err == nil {
		return []byte(v)
	} else {
		response := err.(*gomemcached.MCResponse)
		log.Println(response)
		return nil
	}
}

func (ds *couchbaseDatastore) set(k string, v []byte, ttl int) error {
	if err := ds.validKey(k); err != nil {
		return INVALID_KEY
	}

	if err := ds.validValue(v); err != nil {
		return err
	}

	if ttl > MAX_TTL_IN_SEC {
		ttl = MAX_TTL_IN_SEC
	} else if ttl < 0 {
		ttl = 0
	}

	if err := (*couchbase.Bucket)(ds).SetRaw(k, ttl, v); err != nil {
		response := err.(*gomemcached.MCResponse)
		log.Println(response)
		switch (*response).Status {
		case gomemcached.E2BIG:
			return OVERSIZED_BODY
		default:
			return err
		}
	}
	return nil
}

func (ds *couchbaseDatastore) delete(k string) error {
	if err := ds.validKey(k); err != nil {
		return INVALID_KEY
	}

	if err := (*couchbase.Bucket)(ds).Delete(k); err != nil {
		response := err.(*gomemcached.MCResponse)
		log.Println(response)
		switch (*response).Status {
		case gomemcached.KEY_ENOENT:
			return NOT_FOUND_ERROR
		default:
			return err
		}
	}
	return nil
}

func (ds *couchbaseDatastore) append(k string, v []byte) error {
	if err := ds.validKey(k); err != nil {
		return INVALID_KEY
	}

	if err := ds.validValue(v); err != nil {
		return err
	}

	if err := (*couchbase.Bucket)(ds).Append(k, v); err != nil {
		response := err.(*gomemcached.MCResponse)
		log.Println(response)
		switch (*response).Status {
		case gomemcached.NOT_STORED:
			return NOT_FOUND_ERROR
		case gomemcached.E2BIG:
			return OVERSIZED_BODY
		default:
			return err
		}
	}
	return nil
}

func (ds *couchbaseDatastore) validKey(key string) error {
	if len(key) < 1 || len(key) > MAX_KEY_LENGTH {
		return INVALID_KEY
	}
	return nil
}

func (ds *couchbaseDatastore) validValue(v []byte) error {
	if len(v) == 0 {
		log.Println("body is empty")
		return EMPTY_BODY
	}

	if len(v) > MAX_SIZE_IN_BYTE {
		log.Println("body is too large")
		return OVERSIZED_BODY
	}

	return nil
}
