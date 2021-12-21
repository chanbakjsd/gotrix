package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = DirectEvent{}

// DirectEvent is an event that lists all the DM channels the user is in.
// It is saved in AccountData.
type DirectEvent struct {
	*EventInfo
	Rooms map[matrix.UserID][]matrix.RoomID
}
