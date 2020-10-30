package matrix

import (
	"time"
)

const (
	msInSec = 1000
	nsInMs  = 1000000
)

// Timestamp is a timestamp in milliseconds since Unix time.
type Timestamp int64

// Time converts a timestamp into valid Go time.
func (t Timestamp) Time() time.Time {
	seconds := int64(t / msInSec)
	milliseconds := int64(t % msInSec)
	return time.Unix(seconds, milliseconds*nsInMs)
}
