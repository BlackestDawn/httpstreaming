package request

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

func (r *Request) getContentLength() (int, error) {
	contentLength, err := r.Headers.Get("content-length")
	if err != nil {
		return 0, err
	}
	contentLengthInt, err := strconv.Atoi(contentLength)
	if err != nil {
		return 0, fmt.Errorf("invalid content-length: %s", contentLength)
	}
	return contentLengthInt, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case parseStateInitialized:
		//fmt.Println("DEBUG:: parsing request line")
		requestLine, numBytesRead, err := parseRequestLine(data)
		if numBytesRead != 0 && err == nil {
			r.RequestLine = *requestLine
			r.State = parseStateHeaders
		}
		return numBytesRead, err
	case parseStateHeaders:
		//fmt.Println("DEBUG:: parsing headers")
		numBytesRead, done, err := r.Headers.Parse(data)
		if done {
			r.State = parseStateBody
		}
		return numBytesRead, err
	case parseStateBody:
		//fmt.Println("DEBUG:: parsing body")
		return r.parseBody(data)
	case parseStateDone:
		return 0, fmt.Errorf("request already parsed")
	default:
		return 0, fmt.Errorf("invalid parsing state: %d", r.State)
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// Prep and initialize
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{
		Headers: headers.NewHeaders(),
		Body:    make([]byte, 0),
		State:   parseStateInitialized,
	}

	// Read and loop over data
	for request.State != parseStateDone {
		// Extedn buffer if needed
		if readToIndex >= len(buffer) {
			buffer = append(buffer, make([]byte, bufferSize)...)
		}

		// Read in data
		numBytesReader, err := reader.Read(buffer[readToIndex:])
		readToIndex += numBytesReader

		//fmt.Printf("DEBUG:: data already in buffer '%s'\n", string(buffer))
		//fmt.Printf("DEBUG:: data in last read '%s'\n", string(buffer))
		// Handle errors and EOF
		if err != nil {
			if errors.Is(err, io.EOF) {
				//fmt.Println("DEBUG:: recieved EOF")
				_, err2 := request.parse(buffer[:readToIndex])
				if err2 != nil {
					return nil, err2
				}
				if request.State == parseStateBody && len(request.Body) == 0 {
					request.State = parseStateDone
					break
				}
				if request.State != parseStateDone {
					return nil, fmt.Errorf("incomplete request, in state '%s' when recieved EOF", parseStateNames[request.State])
				}
				break
			}
			return nil, err
		}

		// Parse data
		numBytesParser, err := request.parse(buffer[:readToIndex])
		if numBytesParser == 0 && err == nil {
			continue
		}
		if err != nil {
			return nil, err
		}

		// Remove parsed part and "resetting" buffer
		copy(buffer, buffer[numBytesParser:])
		readToIndex -= numBytesParser
	}

	return request, nil
}
