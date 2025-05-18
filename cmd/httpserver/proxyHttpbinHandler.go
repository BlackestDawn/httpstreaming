package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net/http"
)

func handlerHTTPbin(w *response.Writer, req *request.Request) {
	// Remove unneeded headers
	req.Headers.Del("Content-Length")
	req.Headers.Del("Connection")
	// Add new headers
	req.Headers.Set("Transfer-Encoding", "chunked")
	req.Headers.Set("Trailers", "X-Content-SHA256, X-Content-Length")

	err := w.WriteStatusLine(response.StatusCodeSuccess)
	if err != nil {
		log.Fatalf("Error writing status line in proxy: %v", err)
	}
	err = w.WriteHeaders(req.Headers)
	if err != nil {
		log.Fatalf("Error writing headers in proxy: %v", err)
	}

	log.Printf("establishing proxy connection\n")
	resp, err := http.Get("https://httpbin.org" + req.RequestLine.RequestTarget)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyContent := []byte{}
	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil {
			if errors.Is(err, http.ErrBodyReadAfterClose) || errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("Error reading response body in proxy: %v", err)
		}
		n, err = w.WriteChunkedBody(buf[:n])
		if err != nil {
			log.Fatalf("Error writing response body in proxy: %v", err)
		}
		bodyContent = append(bodyContent, buf[:n]...)
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Fatalf("Error writing chunked body done in proxy: %v", err)
	}

	trailers := headers.NewHeaders()
	checksum := fmt.Sprintf("%x", sha256.Sum256(bodyContent))
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(bodyContent)))
	trailers.Set("X-Content-SHA256", checksum)
	err = w.WriteTrailers(trailers)
	if err != nil {
		log.Fatalf("Error writing trailers in proxy: %v", err)
	}

}
