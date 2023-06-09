package l27

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strconv"
)

func parseContentDispositionFilename(resp *http.Response, fallback string) string {
	contentDisp := resp.Header.Get("Content-Disposition")
	if contentDisp != "" {
		_, params, err := mime.ParseMediaType(contentDisp)
		if err != nil {
			return fallback
		}

		return params["filename"]
	}

	return fallback
}

// Type of all integers returned by the API.
type IntID = int32

// Type of unix time stamps returned by the API.
type IntTime = int64

// In some cases, the API exposes internal int values for statuses instead of string names.
type IntStatus = int32

type ParameterList = BuggyMap[string, ParameterValue]
type ParameterValue = interface{}

// Parse an ID number for the API.
// Returns 0, err if the string is not a valid ID.
func ParseID(id string) (IntID, error) {
	val, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID")
	}

	return IntID(val), nil
}

// Common parameters related to filtering and such that are common to all get-like operations.
type CommonGetParams struct {
	OrderBy   string
	OrderType string
	Filter    string
	PageableParams
}

// Common parameters for APIs that support pagination.
type PageableParams struct {
	Limit  int32
	Offset int32
}

// Workaround for buggy API-responses that return JSON [] when it should be {}.
// Accepts [] and deserializes it as an empty map.
// Task: PL-7612
type BuggyMap[K comparable, V any] struct {
	Map map[K]V
}

func (b *BuggyMap[K, V]) UnmarshalJSON(data []byte) error {
	origErr := json.Unmarshal(data, &b.Map)
	if origErr == nil {
		return nil
	}

	// Try seeing if it reads as an array instead.
	var array []interface{}
	arrayErr := json.Unmarshal(data, &array)
	if arrayErr != nil {
		// Return the map error since that's probably more meaningful.
		return origErr
	}

	if len(array) > 0 {
		return fmt.Errorf("buggy map had an array with more than 0 elements")
	}

	// 0 elements, buggy map.
	b.Map = map[K]V{}
	return nil
}

func (b *BuggyMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Map)
}

// Workaround for API definitions that define [] instead of {} for empty objects.
// Accepts [] and deserializes it as default value.
type EmptyArrayDefault[T any] struct {
	Value T
}

func (b *EmptyArrayDefault[T]) UnmarshalJSON(data []byte) error {
	origErr := json.Unmarshal(data, &b.Value)
	if origErr == nil {
		return nil
	}

	// Try seeing if it reads as an array instead.
	var array []interface{}
	arrayErr := json.Unmarshal(data, &array)
	if arrayErr != nil {
		// Return the original error since that's probably more meaningful.
		return origErr
	}

	if len(array) > 0 {
		return fmt.Errorf("EmptyArrayDefault had an array with more than 0 elements")
	}

	// 0 elements, assume default value
	*b = EmptyArrayDefault[T]{}
	return nil
}

func (b *EmptyArrayDefault[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Value)
}

// A string that may contain a variant for multiple languages, such as "en" and "nl".
type LocalizedString map[string]string

func shallowCloneMap[K comparable, V any](m map[K]V) map[K]V {
	result := map[K]V{}
	for k, v := range m {
		result[k] = v
	}

	return result
}
