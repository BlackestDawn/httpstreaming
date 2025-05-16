package request

import (
	"bytes"
	"fmt"
	"strings"
)

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// Get end-index of line
	idx := bytes.Index(data, []byte(lineBreak))
	if idx == -1 {
		return nil, 0, nil
	}

	// Extract and split into parts
	line := string(data[:idx])
	parts := strings.Split(line, " ")

	// Check request line has 3 parts
	if len(parts) != 3 {
		return nil, idx, fmt.Errorf("invalid request line: %s", line)
	}

	// Validate method
	for _, c := range parts[0] {
		if c < 'A' || c > 'Z' {
			return nil, idx, fmt.Errorf("invalid method: %s", parts[0])
		}
	}

	// Validate HTTP-version
	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, idx, fmt.Errorf("invalid HTTP-version: %s", parts[2])
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}, idx + 2, nil
}
