package json

import "testing"

func TestSpecExamples(t *testing.T) {
	for _, v := range testCases {
		got := string(Canonical([]byte(v.Original)))

		if got != v.Canonical {
			t.Errorf("spec example mismatch:\n"+
				"Got:      %s\n"+
				"Expected: %s",
				got, v.Canonical,
			)
		}
	}
}

func TestUnsignedAndSignaturesDropped(t *testing.T) {
	result := string(Canonical([]byte(`{"unsigned":{"a":"b"}, "signatures":{"c":"d"}}`)))
	if result != "{}" {
		t.Errorf("Unsigned/Signature drop mismatch:\n"+
			"Got:      %s\n"+
			"Expected: {}",
			result,
		)
	}
}

var testCases = []struct {
	Original  string
	Canonical string
}{
	{`{}`, `{}`},
	{`{"b":"2","a":"1"}`, `{"a":"1","b":"2"}`},
	{`{
			"one": 1,
			"two": "Two"
		}`, `{"one":1,"two":"Two"}`},
	{`{
			"b": "2",
			"a": "1"
		}`, `{"a":"1","b":"2"}`},
	{`{
		"auth": {
			"success": true,
			"mxid": "@john.doe:example.com",
			"profile": {
				"display_name": "John Doe",
				"three_pids": [
					{
						"medium": "email",
						"address": "john.doe@example.org"
					},
					{
						"medium": "msisdn",
						"address": "123456789"
					}
				]
			}
		}
	}`, `{"auth":{"mxid":"@john.doe:example.com","profile":{"display_name":"John Doe","three_pids":[{"address":"john.doe@example.org","medium":"email"},{"address":"123456789","medium":"msisdn"}]},"success":true}}`},
	{`{
		"a": "日本語"
	}`, `{"a":"日本語"}`},
	{`{
		"本": 2,
		"日": 1
	}`, `{"日":1,"本":2}`},
	{`{
		"a": "\u65E5"
	}`, `{"a":"日"}`},
	{`{
		"a": null
	}`, `{"a":null}`},
}
