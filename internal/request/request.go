package request

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"log"
	"strconv"
)

func (r *Request) parse(data []byte) (int, error) {
	bytesParsed := 0
	for r.State != parseStateDone {
		numBytesParsed, err := r.parseSingle(data[bytesParsed:])
		if err != nil {
			return 0, err
		}
		if numBytesParsed == 0 {
			break
		}
		bytesParsed += numBytesParsed
	}
	return bytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case parseStateInitialized:
		requestLine, numBytesRead, err := parseRequestLine(data)
		if numBytesRead != 0 && err == nil {
			r.RequestLine = *requestLine
			r.State = parseStateHeaders
		}
		return numBytesRead, err
	case parseStateHeaders:
		numBytesRead, done, err := r.Headers.Parse(data)
		if done {
			r.State = parseStateBody
		}
		return numBytesRead, err
	case parseStateBody:
		// Get content length
		contentLength, exist := r.Headers.Get("Content-Length")
		if !exist {
			// Assume 0 body if content-length header isn't present
			r.State = parseStateDone
			return len(data), nil
		}

		contentLengthInt, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("invalid content-length: %s", contentLength)
		}

		// Append data
		r.Body = append(r.Body, data...)
		if len(r.Body) > contentLengthInt {
			return 0, fmt.Errorf("request body too large")
		}
		if len(r.Body) == contentLengthInt {
			r.State = parseStateDone
		}

		return len(data), nil
	case parseStateDone:
		log.Printf("data left when already parsed: %v\n", data)
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
		if request.State == parseStateDone {
			log.Println("new iteration with done state")
		}

		// Extedn buffer if needed
		if readToIndex >= len(buffer) {
			buffer = append(buffer, make([]byte, bufferSize)...)
		}

		// Read in data
		numBytesReader, err := reader.Read(buffer[readToIndex:])
		readToIndex += numBytesReader

		// Handle errors and EOF
		if err != nil {
			if errors.Is(err, io.EOF) {
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
