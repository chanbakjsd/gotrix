package encrypt

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/event"
)

func init() {
	event.RegisterDefault(TypeRoomEncryption, parseRoomEncryption)
	event.RegisterDefault(TypeRoomEncrypted, parseRoomEncrypted)
	event.RegisterDefault(TypeRoomKey, parseRoomKey)
	event.RegisterDefault(TypeRoomKeyRequest, parseRoomKeyRequest)
	event.RegisterDefault(TypeForwardedRoomKey, parseForwardRoomKey)
	event.RegisterDefault(TypeDummy, parseDummy)
	event.RegisterDefault(TypeRoomKeyWithheld, parseRoomKeyWithheld)
}

func parseRoomEncryption(content json.RawMessage) (event.Event, error) {
	var v RoomEncryptionEvent
	err := json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseRoomEncrypted(content json.RawMessage) (event.Event, error) {
	var v RoomEncryptedEvent
	err := json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseRoomKey(content json.RawMessage) (event.Event, error) {
	var v RoomKeyEvent
	err := json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseRoomKeyRequest(content json.RawMessage) (event.Event, error) {
	var v RoomKeyRequestEvent
	err := json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseForwardRoomKey(content json.RawMessage) (event.Event, error) {
	var v ForwardedRoomKeyEvent
	err := json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseDummy(content json.RawMessage) (event.Event, error) {
	var v DummyEvent
	err := json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseRoomKeyWithheld(content json.RawMessage) (event.Event, error) {
	var v RoomKeyWithheldEvent
	err := json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
