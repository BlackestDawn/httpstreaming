package request

import (
	"fmt"
)

func (r *Request) parseBody(data []byte) (int, error) {
	// Get content length
	contentLength, err := r.getContentLength()
	if err != nil {
		r.State = parseStateDone
		return 0, nil
	}

	// Append data
	r.Body = append(r.Body, data...)
	if len(r.Body) > contentLength {
		return len(data), fmt.Errorf("request body too large")
	}
	if len(r.Body) == contentLength {
		//fmt.Printf("DEBUG:: marking body parsing done\n")
		r.State = parseStateDone
	}

	return len(data), nil
}
