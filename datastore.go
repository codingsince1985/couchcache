package main

import (
	"errors"
)

type datastorer interface {
	get(k string) []byte
	set(k string, v []byte, ttl int) error
	delete(k string) error
	append(k string, v []byte) error
	invalid(k string) bool
}

var (
	NOT_FOUND_ERROR = errors.New("NOT_FOUND")
	TOO_BIG_ERROR   = errors.New("TOO_BIG")
	EMPTY_BODY      = errors.New("EMPTY_BODY")
	INVALID_KEY     = errors.New("INVALID_KEY")
)
