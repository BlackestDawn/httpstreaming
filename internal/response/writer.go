package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"net"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeader
	writerStateBody
	writerStateTrailers
	writerStateClose
)

type writerStateNames map[writerState]string

var writerStateName = writerStateNames{
	writerStateStatusLine: "status line",
	writerStateHeader:     "header",
	writerStateBody:       "body",
	writerStateTrailers:   "trailers",
	writerStateClose:      "close",
}

func (s writerState) String() string {
	return writerStateName[s]
}

type Writer struct {
	writer net.Conn
	state  writerState
}

func NewWriter(w net.Conn) *Writer {
	return &Writer{
		writer: w,
		state:  writerStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) (err error) {
	// Check if at proper state
	if w.state != writerStateStatusLine {
		return fmt.Errorf("invalid state writing status line: %s", w.state)
	}

	if _, ok := statusCodeText[statusCode]; ok {
		_, err = fmt.Fprintf(w.writer, "HTTP/1.1 %d %s\r\n", statusCode, statusCode)
	} else {
		_, err = fmt.Fprintf(w.writer, "HTTP/1.1 %d\r\n", statusCode)
	}
	w.state = writerStateHeader
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	// Check if at proper state
	if w.state != writerStateHeader {
		return fmt.Errorf("invalid state writing headers: %s", w.state)
	}

	for k, v := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w.writer, "\r\n")
	w.state = writerStateBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	// Check if at proper state
	if w.state != writerStateBody {
		return 0, fmt.Errorf("invalid state writing body: %s", w.state)
	}

	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	w.state = writerStateTrailers
	return n, nil
}

func (w *Writer) WriteTrailers(headers headers.Headers) error {
	// Check if at proper state
	if w.state != writerStateTrailers {
		return fmt.Errorf("invalid state writing trailers: %s", w.state)
	}

	for k, v := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w.writer, "\r\n")
	return err
}
