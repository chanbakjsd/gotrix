package encrypt

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

// File represents an encrypted file.
type File struct {
	URL        matrix.URL        `json:"url"`
	Key        JSONWebKey        `json:"key"`
	InitVector string            `json:"iv"`
	Hashes     map[string]string `json:"hashes"`
	Version    string            `json:"v"`
}

// TODO Add helper functions for encryption/decryption.
