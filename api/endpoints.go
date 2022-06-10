package api

import (
	"net/url"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// SupportedVersions contains versions supported by gotrix.
var SupportedVersions = map[string]string{
	"r0.6.0": "r0",
	"r0.6.1": "r0",
	"v1.1":   "v3",
	"v1.2":   "v3",
}

const EndpointSupportedVersions = "_matrix/client/versions"

// Endpoints contains known Matrix Client-Server API endpoints.
type Endpoints struct {
	// Version is the Matrix version for these endpoints.
	Version string
}

func (e Endpoints) Base() string { return "_matrix/client/" + e.Version }

func (e Endpoints) User(id matrix.UserID) string {
	return e.Base() + "/user/" + url.PathEscape(string(id))
}
func (e Endpoints) Room(id matrix.RoomID) string {
	return e.Base() + "/rooms/" + url.PathEscape(string(id))
}
func (e Endpoints) UserRoom(userID matrix.UserID, roomID matrix.RoomID) string {
	return e.User(userID) + "/rooms/" + url.PathEscape(string(roomID))
}

func (e Endpoints) Login() string     { return e.Base() + "/login" }
func (e Endpoints) Logout() string    { return e.Base() + "/logout" }
func (e Endpoints) LogoutAll() string { return e.Logout() + "/all" }

func (e Endpoints) Register() string          { return e.Base() + "/register" }
func (e Endpoints) RegisterAvailable() string { return e.Register() + "/available" }
func (e Endpoints) RegisterRequestToken(authType string) string {
	return e.Register() + "/" + authType + "/requestToken"
}

func (e Endpoints) Account() string           { return e.Base() + "/account" }
func (e Endpoints) AccountWhoami() string     { return e.Account() + "/whoami" }
func (e Endpoints) AccountDeactivate() string { return e.Account() + "/deactivate" }
func (e Endpoints) AccountPassword() string   { return e.Account() + "/password" }
func (e Endpoints) AccountPasswordRequestToken(authType string) string {
	return e.AccountPassword() + "/" + authType + "/requestToken"
}

func (e Endpoints) Account3PID() string     { return e.Account() + "/3pid" }
func (e Endpoints) Account3PIDAdd() string  { return e.Account3PID() + "/add" }
func (e Endpoints) Account3PIDBind() string { return e.Account3PID() + "/bind" }
func (e Endpoints) Account3PIDRequestToken(authType string) string {
	return e.Account3PID() + "/" + authType + "/requestToken"
}

func (e Endpoints) Capabilities() string { return e.Base() + "/capabilities" }
func (e Endpoints) JoinedRooms() string  { return e.Base() + "/joined_rooms" }
func (e Endpoints) PublicRooms() string  { return e.Base() + "/publicRooms" }
func (e Endpoints) Sync() string         { return e.Base() + "/sync" }

func (e Endpoints) Filter(userID matrix.UserID) string {
	return e.User(userID) + "/filter"
}
func (e Endpoints) FilterGet(userID matrix.UserID, filterID string) string {
	return e.Filter(userID) + "/" + url.PathEscape(filterID)
}

func (e Endpoints) RoomCreate() string                      { return e.Base() + "/createRoom" }
func (e Endpoints) RoomAliases(roomID matrix.RoomID) string { return e.Room(roomID) + "/aliases" }
func (e Endpoints) RoomBan(roomID matrix.RoomID) string     { return e.Room(roomID) + "/ban" }
func (e Endpoints) RoomForget(roomID matrix.RoomID) string  { return e.Room(roomID) + "/forget" }
func (e Endpoints) RoomInvite(roomID matrix.RoomID) string  { return e.Room(roomID) + "/invite" }
func (e Endpoints) RoomJoin(roomID matrix.RoomID) string    { return e.Room(roomID) + "/join" }
func (e Endpoints) RoomJoinedMembers(roomID matrix.RoomID) string {
	return e.Room(roomID) + "/joined_members"
}
func (e Endpoints) RoomKick(roomID matrix.RoomID) string    { return e.Room(roomID) + "/kick" }
func (e Endpoints) RoomLeave(roomID matrix.RoomID) string   { return e.Room(roomID) + "/leave" }
func (e Endpoints) RoomMembers(roomID matrix.RoomID) string { return e.Room(roomID) + "/members" }
func (e Endpoints) RoomMessages(roomID matrix.RoomID) string {
	return e.Room(roomID) + "/messages"
}
func (e Endpoints) RoomState(roomID matrix.RoomID) string   { return e.Room(roomID) + "/state" }
func (e Endpoints) RoomUnban(roomID matrix.RoomID) string   { return e.Room(roomID) + "/unban" }
func (e Endpoints) RoomUpgrade(roomID matrix.RoomID) string { return e.Room(roomID) + "/upgrade" }
func (e Endpoints) RoomEvent(roomID matrix.RoomID, eventID matrix.EventID) string {
	return e.Room(roomID) + "/event/" + url.PathEscape(string(eventID))
}
func (e Endpoints) RoomStateExact(roomID matrix.RoomID, eventType event.Type, stateKey string) string {
	return e.RoomState(roomID) + "/" + url.PathEscape(string(eventType)) + "/" + url.PathEscape(stateKey)
}
func (e Endpoints) RoomReceipt(roomID matrix.RoomID, receiptType ReceiptType, eventID matrix.EventID) string {
	return e.RoomState(roomID) + "/" + url.PathEscape(string(receiptType)) + "/" + url.PathEscape(string(eventID))
}
func (e Endpoints) RoomRedact(roomID matrix.RoomID, eventID matrix.EventID, transactionID string) string {
	return e.Room(roomID) + "/redact/" + url.PathEscape(string(eventID)) + "/" + url.PathEscape(transactionID)
}
func (e Endpoints) RoomSend(roomID matrix.RoomID, eventType event.Type, transactionID string) string {
	return e.Room(roomID) + "/send/" + url.PathEscape(string(eventType)) + "/" + url.PathEscape(transactionID)
}
func (e Endpoints) RoomTyping(roomID matrix.RoomID, userID matrix.UserID) string {
	return e.Room(roomID) + "/typing/" + url.PathEscape(string(userID))
}

func (e Endpoints) Directory() string { return e.Base() + "/directory" }
func (e Endpoints) DirectoryRoomAlias(roomAlias string) string {
	return e.Directory() + "/room/" + url.PathEscape(roomAlias)
}
func (e Endpoints) DirectoryListRoom(roomID matrix.RoomID) string {
	return e.Directory() + "/list/room/" + url.PathEscape(string(roomID))
}

func (e Endpoints) UserDirectorySearch() string { return e.Base() + "/user_directory/search" }

func (e Endpoints) Profile(userID matrix.UserID) string {
	return e.Base() + "/profile/" + url.PathEscape(string(userID))
}
func (e Endpoints) ProfileAvatarURL(userID matrix.UserID) string {
	return e.Profile(userID) + "/avatar_url"
}
func (e Endpoints) ProfileDisplayName(userID matrix.UserID) string {
	return e.Profile(userID) + "/displayname"
}

func (e Endpoints) Media() string           { return "_matrix/media/" + e.Version }
func (e Endpoints) MediaConfig() string     { return e.Media() + "/config" }
func (e Endpoints) MediaPreviewURL() string { return e.Media() + "/preview_url" }
func (e Endpoints) MediaUpload() string     { return e.Media() + "/upload" }
func (e Endpoints) MediaDownload(serverName string, mediaID string, fileName string) string {
	return e.Media() + "/download/" + url.PathEscape(serverName) + "/" + url.PathEscape(mediaID) + "/" +
		url.PathEscape(fileName)
}
func (e Endpoints) MediaThumbnail(serverName string, mediaID string) string {
	return e.Media() + "/thumbnail/" + url.PathEscape(serverName) + "/" + url.PathEscape(mediaID)
}

func (e Endpoints) VOIPTURNServers() string { return e.Media() + "/voip/turnServer" }

func (e Endpoints) PresenceStatus(userID matrix.UserID) string {
	return e.Base() + "/presence/" + url.PathEscape(string(userID)) + "/status"
}

func (e Endpoints) AccountDataGlobal(userID matrix.UserID, dataType string) string {
	return e.User(userID) + "/account_data/" + url.PathEscape(dataType)
}
func (e Endpoints) AccountDataRoom(userID matrix.UserID, roomID matrix.RoomID, dataType string) string {
	return e.UserRoom(userID, roomID) + "/account_data/" + url.PathEscape(dataType)
}

func (e Endpoints) SendToDevice(eventType event.Type, transactionID string) string {
	return e.Base() + "/sendToDevice/" + url.PathEscape(string(eventType)) + "/" + url.PathEscape(transactionID)
}

func (e Endpoints) Tags(userID matrix.UserID, roomID matrix.RoomID) string {
	return e.UserRoom(userID, roomID) + "/tags"
}

func (e Endpoints) Tag(userID matrix.UserID, roomID matrix.RoomID, name matrix.TagName) string {
	return e.Tags(userID, roomID) + "/" + url.PathEscape(string(name))
}

func (e Endpoints) SSOLogin(redirectURL string) string {
	return e.Base() + "/login/sso/redirect?redirectUrl=" + url.QueryEscape(redirectURL)
}
