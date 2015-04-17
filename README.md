Caching solution backed by Couchbase
=

A caching service developed in Go. It provides REST APIs to populate and retrieve key-value pairs stored in Couchbase.

To populate a key-value pair
-
POST request to `http://your-host-name:your-port/key/your-key[?ttl=time-to-live-in-sec]` with data in body

To retrieve a key-value pair
-
GET request to `http://your-host-name:your-port/key/your-key`

License
=
couchcache is distributed under the terms of the MIT license. See LICENSE for details.