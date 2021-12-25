package matrix

import "encoding/json"

// Capabilities represent the server's capabilities as defined in
// https://spec.matrix.org/v1.1/client-server-api/#get_matrixclientv3capabilities.
type Capabilities map[string]json.RawMessage

// CapabilityChangePassword contains whether the server has password change enabled.
// If it is disabled, it's probably delegated to the identity server.
type CapabilityChangePassword struct {
	Enabled bool `json:"enabled"`
}

// CapabilityRoomVersion contains information of the version of the rooms.
type CapabilityRoomVersion struct {
	Default   string `json:"default"`
	Available map[string]RoomVersionStability
}

// RoomVersionStability expresses the stability status of a room version.
type RoomVersionStability string

// RoomVersionStability can either be "stable" or "unstable".
const (
	RoomVersionStable   RoomVersionStability = "stable"
	RoomVersionUnstable RoomVersionStability = "unstable"
)

// ChangePassword retrieves CapabilityChangePassword from the Capabilities.
func (c Capabilities) ChangePassword() (CapabilityChangePassword, error) {
	var resp CapabilityChangePassword
	err := json.Unmarshal(c["m.change_password"], &resp)
	return resp, err
}

// RoomVersion retrieves CapabilityRoomVersion from the Capabilities.
func (c Capabilities) RoomVersion() (CapabilityRoomVersion, error) {
	var resp CapabilityRoomVersion
	err := json.Unmarshal(c["m.room_versions"], &resp)
	return resp, err
}
