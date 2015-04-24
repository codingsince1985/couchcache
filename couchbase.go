package main

import (
	"github.com/couchbase/go-couchbase"
)

const MAX_TTL_IN_SEC = 60 * 60 * 24 * 30

type couchbaseDatastore couchbase.Bucket

func newDatastore(url string, bucket, pass *string) (ds *couchbaseDatastore, err error) {
	c, err := couchbase.ConnectWithAuthCreds(url, *bucket, *pass)
	if err != nil {
		return ds, err
	}

	p, err := c.GetPool("default")
	if err != nil {
		return ds, err
	}

	b, err := p.GetBucketWithAuth(*bucket, *bucket, *pass)
	if err != nil {
		return ds, err
	}

	return (*couchbaseDatastore)(b), nil
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
