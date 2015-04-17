package main

import (
	"github.com/couchbase/go-couchbase"
)

const MAX_TTL_IN_SEC = 2592000

type datastore couchbase.Bucket

func newDatastore(url, bucket string) (ds datastore, err error) {
	c, err := couchbase.Connect(url)
	if err != nil {
		return ds, err
	}

	p, err := c.GetPool("default")
	if err != nil {
		return ds, err
	}

	b, err := p.GetBucketWithAuth(bucket, bucket, "")
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
