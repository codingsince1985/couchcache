Couchcache
=

A caching service developed in Go. It provides REST APIs to access key-value pairs stored in Couchbase.

To start couchcache
-
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
-
`http://HOST:8080/key/KEY`

To store a key-value pair
-
Send `POST` request to endpoint with data in body, optionally set TTL by providing `?ttl=TTL_IN_SEC`

To retrieve a key
-
Send `GET` request to endpoint

To delete a key
-
Send `DELETE` request to endpoint

To append data for a key
-
Send `PUT` request to endpoint with data in body

Limitation
=
* Max key length is 250 bytes
* Max value size is 20 MB

See [Couchbase Limits](http://docs.couchbase.com/admin/admin/Misc/limits.html).

License
=
couchcache is distributed under the terms of the MIT license. See LICENSE for details.
