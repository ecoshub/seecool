package seecool

import "errors"

var (
	errMalformedQuery error = errors.New("Malformed Query. error_code:00")
	errMissingKVQuery error = errors.New("Missing key and value fro INSERT query. error_code:01")
)
