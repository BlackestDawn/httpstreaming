package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

var statusText = map[StatusCode]string{
	StatusCodeSuccess:             "OK",
	StatusCodeBadRequest:          "BAD REQUEST",
	StatusCodeInternalServerError: "INTERNAL SERVER ERROR",
}

func (s StatusCode) String() string {
	return statusText[s]
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	if text, ok := statusText[statusCode]; ok {
		_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, text)
		return err
	}
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d\r\n", statusCode)
	return err
}

func GetDefaultHeaders(contentlen int) headers.Headers {
	return headers.Headers{
		"Content-Length": fmt.Sprintf("%d", contentlen),
		"Content-Type":   "text/plain",
		"Connection":     "close",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\r\n")
	return err
}
