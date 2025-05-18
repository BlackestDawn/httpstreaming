package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	fmt.Fprintf(w.writer, "%x\r\n%s\r\n", len(p), p)
	return len(p), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	fmt.Fprintf(w.writer, "0\r\n")
	return 0, nil
}
