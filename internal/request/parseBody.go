package request

import (
	"fmt"
	"strconv"
)

func (r *Request) parseBody(data []byte) (int, error) {
	// Get content length
	contentLength, exist := r.Headers.Get("Content-Length")
	if !exist {
		// Assume 0 body if content-length header isn't present
		r.State = parseStateDone
		return 0, nil
	}

	contentLengthInt, err := strconv.Atoi(contentLength)
	if err != nil {
		return 0, fmt.Errorf("invalid content-length: %s", contentLength)
	}

	// Append data
	r.Body = append(r.Body, data...)
	if len(r.Body) > contentLengthInt {
		return len(data), fmt.Errorf("request body too large")
	}
	if len(r.Body) == contentLengthInt {
		r.State = parseStateDone
	}

	return len(data), nil
}
