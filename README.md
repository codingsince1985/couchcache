Caching solution backed by Couchbase
=

A caching service developed in Go. It provides REST APIs to populate and retrieve key-value pairs stored in Couchbase.

To start couchcache
-
Run couchcache with Couchbase server information

`./couchcache --host=hostname --port=port --bucket=bucketname --pass=password`

default values
```
hostname = localhost
port = 8091
bucketname = couchcache
pass = password
```

To store a key-value pair
-
Send POST request to `http://your-host-name:your-port/key/your-key[?ttl=time-to-live-in-sec]` with data in body

To retrieve a key-value pair
-
Send GET request to `http://your-host-name:your-port/key/your-key`

To delete a key-value pair
-
Send DELETE request to `http://your-host-name:your-port/key/your-key`

License
=
couchcache is distributed under the terms of the MIT license. See LICENSE for details.