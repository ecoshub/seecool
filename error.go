package seecool

import "errors"

var (
	errMalformedQuery error = errors.New("Malformed Query.")
	errMissingKVQuery error = errors.New("Missing key and value fro INSERT query.")
)
