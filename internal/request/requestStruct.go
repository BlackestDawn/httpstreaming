package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
)

type parsingState int

const (
	parseStateInitialized parsingState = iota
	parseStateHeaders
	parseStateBody
	parseStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	State       parsingState
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

const (
	bufferSize int    = 8
	lineBreak  string = "\r\n"
)

var parseStateNames = map[parsingState]string{
	parseStateInitialized: "Initialized",
	parseStateHeaders:     "Headers",
	parseStateBody:        "Body",
	parseStateDone:        "Done",
}

func (r *Request) String() string {
	output := `Request line:
- Method: %s
- Target: %s
- Version: %s
`
	output = fmt.Sprintf(output, r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
	output += "Headers:\n"
	for k, v := range r.Headers {
		output += fmt.Sprintf("- %s: %s\n", k, v)
	}
	output += "Body:\n"
	output += string(r.Body)
	return output
}
