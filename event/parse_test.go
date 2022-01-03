package event

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseUnsignedData(t *testing.T) {
	const redactReasonData = `{ "type": "foo_bar" }`
	const prevContentData = `{ "field": 1234 }`
	const data = `
		{
			"age": 456,
			"redacted_because": ` + redactReasonData + `,
			"prev_content": ` + prevContentData + `,
			"test_future_field_compliance": 56789,
			"transaction_id": "abcdef\ngh"
		}
	`

	expected := UnsignedData{
		Age:           456,
		RedactReason:  RawEvent(redactReasonData),
		PrevContent:   json.RawMessage(prevContentData),
		TransactionID: "abcdef\ngh",
	}

	var actual UnsignedData
	err := json.Unmarshal([]byte(data), &actual)
	if err != nil {
		t.Fatalf("unexpected error parsing unsigned data: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("mismatch on parsing unsigned data\nexpected: %#v\ngot: %#v", expected, actual)
	}
}

func TestRawEventCopied(t *testing.T) {
	const data = `
			{
				"content": {
					"@bob:example.com": [
						"!abcdefgh:example.com",
						"!hgfedcba:example.com"
					]
				},
				"type": "m.direct"
			}
	`

	a := []byte(data)
	v, err := Parse(RawEvent(a))
	if err != nil {
		t.Fatalf("unexpected error parsing event: %v", err)
	}

	// Corrupt passed in RawEvent.
	a[0] = 'k'

	if string(v.Info().Raw) != data {
		t.Fatalf("mismatch between raw data\nexpected: %s\ngot: %s", data, v.Info().Raw)
	}
}
