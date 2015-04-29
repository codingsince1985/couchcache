package main

import (
	"flag"
	"fmt"
	"github.com/couchbase/go-couchbase"
	"github.com/couchbase/gomemcached"
	"log"
)

const MAX_TTL_IN_SEC = 60 * 60 * 24 * 30
const MAX_SIZE_IN_BYTE = 20 * 1024 * 1024

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
	if ttl > MAX_TTL_IN_SEC {
		ttl = MAX_TTL_IN_SEC
	} else if ttl < 0 {
		ttl = 0
	}

	if len(v) == 0 {
		log.Println(k, "is empty")
		return EMPTY_BODY
	}

	if len(v) > MAX_SIZE_IN_BYTE {
		log.Println(k, "is too big")
		return TOO_BIG_ERROR
	}

	if err := (*couchbase.Bucket)(ds).SetRaw(k, ttl, v); err != nil {
		response := err.(*gomemcached.MCResponse)
		log.Println(response)
		switch (*response).Status {
		case gomemcached.E2BIG:
			return TOO_BIG_ERROR
		default:
			return err
		}
	}
	return nil
}

func (ds *couchbaseDatastore) delete(k string) error {
	if err := (*couchbase.Bucket)(ds).Delete(k); err != nil {
		response := err.(*gomemcached.MCResponse)
		log.Println(response)
		if (*response).Status == gomemcached.KEY_ENOENT {
			return nil
		}
		return err
	}
	return nil
}

func (ds *couchbaseDatastore) append(k string, v []byte) error {
	if len(v) == 0 {
		log.Println(k, "is empty")
		return EMPTY_BODY
	}

	if len(v) > MAX_SIZE_IN_BYTE {
		log.Println(k, "is too big")
		return TOO_BIG_ERROR
	}

	if err := (*couchbase.Bucket)(ds).Append(k, v); err != nil {
		response := err.(*gomemcached.MCResponse)
		log.Println(response)
		switch (*response).Status {
		case gomemcached.NOT_STORED:
			return NOT_FOUND_ERROR
		case gomemcached.E2BIG:
			return TOO_BIG_ERROR
		default:
			return err
		}
	}
	return nil
}
