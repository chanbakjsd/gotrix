package json

import (
	"encoding/json"
)

// CanonicalMarshal marshals the provided input and returns the canonicalized JSON.
func CanonicalMarshal(input interface{}) ([]byte, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return Canonical(b), nil
}

// Canonical assumes that the JSON is valid and canonicalizes the JSON.
func Canonical(input []byte) []byte {
	return sortJSON(compact(input))
}
