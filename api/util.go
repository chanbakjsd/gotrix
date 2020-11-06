package api

import (
	"encoding/base64"
	"math/rand"
	"time"
)

var transactionID = make(chan string)

// NextTransactionID retrieves the next transaction ID to use.
func NextTransactionID() string {
	return <-transactionID
}

func init() {
	go generateTransactionID()
}

// generateTransactionID saturates the transactionID channel with new IDs.
// This minifies the possibility that a transaction ID collision happening.
func generateTransactionID() {
	// localRand is an instance of random so global random is not affected.
	// It is intentionally not cryptographically secure to make it faster.
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		var bytes [16]byte
		_, err := localRand.Read(bytes[:])
		if err != nil {
			panic(err)
		}
		transactionID <- base64.RawURLEncoding.EncodeToString(bytes[:])
	}
}
