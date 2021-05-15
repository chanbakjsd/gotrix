package json

import (
	"sort"

	"github.com/tidwall/gjson"
)

func sortJSON(input []byte) []byte {
	output := make([]byte, 0, len(input))
	result := gjson.ParseBytes(input)

	rawJSON := rawJSONFromResult(result, input)
	return sortJSONValue(result, rawJSON, output)
}

// rawJSONFromResult extracts the raw JSON bytes pointed to by result.
// input must be the json bytes that were used to generate result
func rawJSONFromResult(result gjson.Result, input []byte) []byte {
	// This is lifted from gjson README. Basically, result.Raw is a copy of
	// the bytes we want, but its more efficient to take a slice.
	// If Index is 0 then for some reason we can't extract it from the original
	// JSON bytes.
	if result.Index > 0 {
		return input[result.Index : result.Index+len(result.Raw)]
	}
	return []byte(result.Raw)
}

// sortJSONValue takes a gjson.Result and sorts it. inputJSON must be the
// raw JSON bytes that gjson.Result points to.
func sortJSONValue(input gjson.Result, inputJSON, output []byte) []byte {
	if input.IsArray() {
		return sortJSONArray(input, inputJSON, output)
	}

	if input.IsObject() {
		return sortJSONObject(input, inputJSON, output)
	}

	// If its neither an object nor an array then there is no sub structure
	// to sort, so just append the raw bytes.
	return append(output, inputJSON...)
}

// sortJSONArray takes a gjson.Result and sorts it, assuming its an array.
// inputJSON must be the raw JSON bytes that gjson.Result points to.
func sortJSONArray(input gjson.Result, inputJSON, output []byte) []byte {
	sep := byte('[')

	// Iterate over each value in the array and sort it.
	input.ForEach(func(_, value gjson.Result) bool {
		output = append(output, sep)
		sep = ','

		rawJSON := rawJSONFromResult(value, inputJSON)
		output = sortJSONValue(value, rawJSON, output)

		return true // keep iterating
	})

	if sep == '[' {
		// If sep is still '[' then the array was empty and we never wrote the
		// initial '[', so we write it now along with the closing ']'.
		output = append(output, '[', ']')
	} else {
		// Otherwise we end the array by writing a single ']'
		output = append(output, ']')
	}
	return output
}

// sortJSONObject takes a gjson.Result and sorts it, assuming its an object.
// inputJSON must be the raw JSON bytes that gjson.Result points to.
func sortJSONObject(input gjson.Result, inputJSON, output []byte) []byte {
	type entry struct {
		key    string // The parsed key string
		rawKey []byte // The raw, unparsed key JSON string
		value  gjson.Result
	}

	var entries []entry

	// Iterate over each key/value pair and add it to a slice
	// that we can sort
	input.ForEach(func(key, value gjson.Result) bool {
		if key.String() == "unsigned" || key.String() == "signatures" {
			// Don't include unsigned/signatures.
			return true
		}
		entries = append(entries, entry{
			key:    key.String(),
			rawKey: rawJSONFromResult(key, inputJSON),
			value:  value,
		})
		return true // keep iterating
	})

	// Sort the slice based on the *parsed* key
	sort.Slice(entries, func(a, b int) bool {
		return entries[a].key < entries[b].key
	})

	sep := byte('{')

	for _, entry := range entries {
		output = append(output, sep)
		sep = ','

		// Append the raw unparsed JSON key, *not* the parsed key
		output = append(output, entry.rawKey...)
		output = append(output, ':')

		rawJSON := rawJSONFromResult(entry.value, inputJSON)

		output = sortJSONValue(entry.value, rawJSON, output)
	}
	if sep == '{' {
		// If sep is still '{' then the object was empty and we never wrote the
		// initial '{', so we write it now along with the closing '}'.
		output = append(output, '{', '}')
	} else {
		// Otherwise we end the object by writing a single '}'
		output = append(output, '}')
	}
	return output
}
