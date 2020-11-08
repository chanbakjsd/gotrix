package encrypt

// JSONWebKey represents a JSON web key.
type JSONWebKey struct {
	KeyType       string   `json:"kty"`     // Type of key. Must be "oct".
	KeyOperations []string `json:"key_ops"` // Key operations.
	Algorithm     string   `json:"alg"`     // URLSafe unpadded base64 key.
	Extractable   bool     `json:"ext"`     // Must be true.
}
