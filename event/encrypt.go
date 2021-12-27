package event

import "github.com/chanbakjsd/gotrix/matrix"

// EncryptedFile represents an encrypted file.
type EncryptedFile struct {
	URL        matrix.URL        `json:"url"`
	Key        JSONWebKey        `json:"key"`
	InitVector string            `json:"iv"`
	Hashes     map[string]string `json:"hashes"`
	Version    string            `json:"v"`
}

// JSONWebKey represents a JSON web key.
type JSONWebKey struct {
	KeyType       string   `json:"kty"`     // Type of key. Must be "oct".
	KeyOperations []string `json:"key_ops"` // Key operations.
	Algorithm     string   `json:"alg"`     // URLSafe unpadded base64 key.
	Extractable   bool     `json:"ext"`     // Must be true.
}
