Couchcache
==

A caching service developed in Go. It provides REST APIs to access key-value pairs stored in Couchbase.

To start couchcache
--
Run couchcache with Couchbase server (host and port) and bucket (name and password) information

`./couchcache --host=HOST --port=PORT --bucket=BUCKET --pass=PASS`

default values
```
host: localhost
port: 8091
bucket: couchcache
pass: password
```
For example, `./couchcache --host=10.99.107.192 --port=8091 --bucket=cachebucket --pass=c@che1t`

Cache service endpoint
--
`http://HOST:8080/key/KEY`

### To store a key-value pair
* request
  * send `POST` request to endpoint with data in body
  * optionally set TTL by `?ttl=TTL_IN_SEC`
* response
  * HTTP 201 if stored
  * HTTP 400 if key or value is invalid

### To retrieve a key
* request
  * send `GET` request to endpoint
* response
  * HTTP 200 with data in body
  * HTTP 404 if key doesn't exist
  * HTTP 400 if key is invalid

### To delete a key
* request
  * send `DELETE` request to endpoint
* response
  * HTTP 200 if deleted
  * HTTP 404 is key doesn't exist
  * HTTP 400 if key is invalid

### To append data for a key
* request
  * send `PUT` request to endpoint with data in body
* response
  * HTTP 200 if appended
  * HTTP 404 if key doesn't exist
  * HTTP 400 if key or value is invalid

Limitation
--
* Max key length is 250 bytes
* Max value size is 20 MB

See [Couchbase Limits](http://docs.couchbase.com/admin/admin/Misc/limits.html).

License
==
couchcache is distributed under the terms of the MIT license. See LICENSE for details.
