package main

import (
	"flag"
	"fmt"
	"gopkg.in/couchbase/gocb.v1"
	"log"
)

const (
	maxTTLInSec   = 60 * 60 * 24 * 30
	maxSizeInByte = 20 * 1024 * 1024
	maxKeyLength  = 250
)

type couchbaseDatastore gocb.Bucket

func newDatastore() (ds *couchbaseDatastore, err error) {
	url, bucket, pass := parseFlag()

	if c, err := gocb.Connect(url); err == nil {
		if b, err := c.OpenBucket(bucket, pass); err == nil {
			return (*couchbaseDatastore)(b), nil
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
	var val []uint8
	if _, err := (*gocb.Bucket)(ds).Get(k, &val); err != nil {
		if err.Error() != "Key not found." {
			log.Println(err)
		}
		return nil
	}
	return []byte(val)
}

func (ds *couchbaseDatastore) set(k string, v []byte, ttl int) error {
	if ttl > maxTTLInSec {
		ttl = maxTTLInSec
	} else if ttl < 0 {
		ttl = 0
	}

	_, err := (*gocb.Bucket)(ds).Upsert(k, v, uint32(ttl))
	return memdErrorToDatastoreError(err)

}

func (ds *couchbaseDatastore) delete(k string) error {
	if err := ds.validKey(k); err != nil {
		return errInvalidKey
	}

	_, err := (*gocb.Bucket)(ds).Remove(k, gocb.Cas(0))
	return memdErrorToDatastoreError(err)
}

func (ds *couchbaseDatastore) append(k string, v []byte) error {
	if err := ds.validKey(k); err != nil {
		return errInvalidKey
	}

	if err := ds.validValue(v); err != nil {
		return err
	}

	_, err := (*gocb.Bucket)(ds).Append(k, string(v))
	return memdErrorToDatastoreError(err)

}

func (ds *couchbaseDatastore) validKey(key string) error {
	if len(key) < 1 || len(key) > maxKeyLength {
		return errInvalidKey
	}
	return nil
}

func (ds *couchbaseDatastore) validValue(v []byte) error {
	if len(v) == 0 {
		log.Println("body is empty")
		return errEmptyBody
	}

	if len(v) > maxSizeInByte {
		log.Println("body is too large")
		return errOversizedBody
	}

	return nil
}

func memdErrorToDatastoreError(err error) error {
	if err == nil {
		return nil
	}

	log.Println(err.Error())
	switch err.Error() {
	case "Key not found.":
		return errNotFound
	case "The document could not be stored.":
		return errNotFound
	case "Document value was too large.":
		return errOversizedBody
	default:
		log.Println(err)
		return err
	}
}
