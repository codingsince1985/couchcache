package main

import (
	"github.com/couchbase/go-couchbase"
)

type pool struct {
	url string
	p   couchbase.Pool
}

func newPool(url string) *pool {
	client, _ := couchbase.Connect(url)
	p, _ := client.GetPool("default")
	return &pool{url, p}
}

func (pool *pool) getBucket(bucket string) (*couchbase.Bucket, error) {
	b, err := pool.p.GetBucketWithAuth(bucket, bucket, "")
	return b, err
}
