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
	errNotFound      = errors.New("NOT_FOUND")
	errKeyExists     = errors.New("KEY_EXISTS_ERROR")
	errOversizedBody = errors.New("OVERSIZED_BODY")
	errEmptyBody     = errors.New("EMPTY_BODY")
	errInvalidKey    = errors.New("INVALID_KEY")
	errInvalidBody   = errors.New("INVALID_BODY")
)
