Caching solution backed by Couchbase
=

A caching service developed in Go. It provides REST APIs to populate and retrieve key-value pairs stored in Couchbase.

To start couchcache
-
Run couchcache with Couchbase server and bucket information

`./couchcache --host=hostname --port=port --bucket=bucketname --pass=password`

default values
```
hostname = localhost
port = 8091
bucketname = couchcache
pass = password
```
Cache service endpoint
-
`http://HOSTNAME:PORT/key/KEY`

To store a key-value pair
-
Send POST request to endpoint with data in body, optionally set TTL by providing `?ttl=TTL_IN_SEC`

To retrieve a key
-
Send GET request to endpoint

To delete a key
-
Send DELETE request to endpoint

To append data for a key
-
Send PUT request to endpoint with data in body

Limitation
=
* Max key length is 250 bytes
* Max value size is 20 MB

See [Couchbase Limits](http://docs.couchbase.com/admin/admin/Misc/limits.html).

License
=
couchcache is distributed under the terms of the MIT license. See LICENSE for details.
