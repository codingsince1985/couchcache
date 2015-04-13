package main

import (
	"github.com/couchbase/go-couchbase"
)

type datastore couchbase.Bucket

func newDatastore(url, bucket string) datastore {
	c, _ := couchbase.Connect(url)
	p, _ := c.GetPool("default")
	b, _ := p.GetBucketWithAuth(bucket, bucket, "")
	return datastore(*b)
}

func (ds *datastore) get(key string) []byte {
	v, _ := (*couchbase.Bucket)(ds).GetRaw(key)
	return []byte(v)
}
