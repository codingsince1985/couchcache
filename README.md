Couchcache
==
[![Go Report Card](https://goreportcard.com/badge/codingsince1985/couchcache)](https://goreportcard.com/report/codingsince1985/couchcache)

A caching service developed in Go. It provides REST APIs to access key-value pairs stored in Couchbase.

You may also consider [using couchcache as a mocked service](http://codingsince1985.blogspot.com.au/2015/05/use-caching-service-as-mocked.html) when doing TDD.

To start couchcache
--
Run couchcache with Couchbase server (host and port) and bucket (name and password) information

`./couchcache --host=HOST --port=PORT --bucket=BUCKET --pass=PASS`
#### Example
`./couchcache --host=10.99.107.192 --port=8091 --bucket=cachebucket --pass=c@che1t`
#### Default values
```
host: localhost
port: 8091
bucket: couchcache
pass: password
```
Cache service endpoint
--
`http://HOST:8080/key/KEY`
#### Examples
`http://10.99.107.190:8080/key/customer_555`

`http://10.99.107.190:8080/key/the_service_i_want_to_mock-endpoint_a`, if you're mocking other service's endpoint

### To store a key-value pair
* request
  * send `POST` request to endpoint with data in body
  * optionally set TTL by `?ttl=TTL_IN_SEC`
* response
  * `HTTP 201 Created` if stored
  * `HTTP 400 Bad Request` if key or value is invalid

### To retrieve a key
* request
  * send `GET` request to endpoint
* response
  * `HTTP 200 OK` with data in body
  * `HTTP 404 Not Found` if key doesn't exist
  * `HTTP 400 Bad Request` if key is invalid

### To delete a key
* request
  * send `DELETE` request to endpoint
* response
  * `HTTP 204 No Content` if deleted
  * `HTTP 404 Not Found` is key doesn't exist
  * `HTTP 400 Bad Request` if key is invalid

### To append data for a key
* request
  * send `PUT` request to endpoint with data in body
* response
  * `HTTP 200 OK` if appended
  * `HTTP 404 Not Found` if key doesn't exist
  * `HTTP 400 Bad Request` if key or value is invalid

Limitations
--
* Max key length is 250 bytes
* Max value size is 20 MB

See [Couchbase Limits](http://docs.couchbase.com/admin/admin/Misc/limits.html).

License
==
couchcache is distributed under the terms of the MIT license. See LICENSE for details.
