package l27

import (
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
