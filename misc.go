package l27

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
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

// Common parameters related to filtering and such that are common to all get-like operations.
type CommonGetParams struct {
	OrderBy   string
	OrderType string
	Limit     int
	Offset    int
	Filter    string
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
