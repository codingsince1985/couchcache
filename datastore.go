package main

import (
	"github.com/couchbase/go-couchbase"
)

const MAX_TTL_IN_SEC = 60 * 60 * 24 * 30

type datastore couchbase.Bucket

func newDatastore(url string, bucket, pass *string) (ds datastore, err error) {
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

	return datastore(*b), nil
}

func (ds *datastore) get(k string) []byte {
	v, err := (*couchbase.Bucket)(ds).GetRaw(k)
	if err != nil {
		return nil
	}

	return []byte(v)
}

func (ds *datastore) set(k string, v []byte, ttl int) error {
	if ttl > MAX_TTL_IN_SEC {
		ttl = MAX_TTL_IN_SEC
	}

	err := (*couchbase.Bucket)(ds).SetRaw(k, ttl, v)
	if err != nil {
		return err
	}

	return nil
}
