package main

import (
	"errors"
)

type datastorer interface {
	get(k string) []byte
	set(k string, v []byte, ttl int) error
	delete(k string) error
	append(k string, v []byte) error
	validKey(k string) error
	validValue(v []byte) error
}

var (
	NOT_FOUND_ERROR  = errors.New("NOT_FOUND")
	KEY_EXISTS_ERROR = errors.New("KEY_EXISTS_ERROR")
	OVERSIZED_BODY   = errors.New("OVERSIZED_BODY")
	EMPTY_BODY       = errors.New("EMPTY_BODY")
	INVALID_KEY      = errors.New("INVALID_KEY")
	INVALID_BODY     = errors.New("INVALID_BODY")
)
