package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const lineBreak = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Get end-index of line
	idx := bytes.Index(data, []byte(lineBreak))
	if idx == -1 {
		return 0, false, nil
	}

	// Check if at end of headers
	if idx == 0 {
		return 2, true, nil
	}

	// Extract and split into parts
	parts := strings.SplitN(string(data[:idx]), ":", 2)

	// Validate key
	key := strings.TrimLeft(parts[0], " ")
	if key != strings.TrimRight(key, " ") || !validateHeaderKey(key) {
		return 0, false, fmt.Errorf("invalid header key: %s", key)
	}

	// Set value
	h.Set(key, parts[1])
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	// Formatting
	key = strings.ToLower(key)
	value = strings.TrimSpace(value)

	// Check if already existing
	if val, ok := h[key]; ok {
		h[key] = val + ", " + value
	} else {
		h[key] = value
	}
}

func validateHeaderKey(key string) bool {
	if len(key) == 0 {
		return false
	}

	// Check for invalid characters
	invalidChars := "(),/:;<=>?@[\\]{}"
	for _, b := range key {
		if b < 33 || b > 126 || strings.ContainsRune(invalidChars, b) {
			return false
		}
	}
	return true
}

func (h Headers) Get(key string) (string, bool) {
	v, ok := h[strings.ToLower(key)]
	return v, ok
}

func (h Headers) Del(key string) {
	delete(h, strings.ToLower(key))
}
