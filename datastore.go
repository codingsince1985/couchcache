package main

import (
	"errors"
)

type datastorer interface {
	get(k string) []byte
	set(k string, v []byte, ttl int) error
	delete(k string) error
	append(k string, v []byte) error
}

var NOT_FOUND_ERROR = errors.New("NOT_FOUND")
var TOO_BIG_ERROR = errors.New("TOO_BIG")
var EMPTY_BODY = errors.New("EMPTY_BODY")
