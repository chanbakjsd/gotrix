package api

import (
	"encoding/base64"
	"math/rand"
	"sync"
	"time"
)

var (
	localRand   = rand.New(rand.NewSource(time.Now().UnixNano()))
	localRandMu sync.Mutex
)

// NextTransactionID retrieves the next transaction ID to use.
func NextTransactionID() string {
	var bytes [16]byte

	localRandMu.Lock()
	_, err := localRand.Read(bytes[:])
	localRandMu.Unlock()

	if err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes[:])
}
