package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

var statusCodeText = map[StatusCode]string{
	StatusCodeSuccess:             "OK",
	StatusCodeBadRequest:          "BAD REQUEST",
	StatusCodeInternalServerError: "INTERNAL SERVER ERROR",
}

func (s StatusCode) String() string {
	return statusCodeText[s]
}

func GetDefaultHeaders(contentlen int) headers.Headers {
	return headers.Headers{
		"Content-Length": fmt.Sprintf("%d", contentlen),
		"Content-Type":   "text/plain",
		"Connection":     "close",
	}
}
