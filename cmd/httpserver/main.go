package main

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

type pageData struct {
	pageTitle   string
	pageHeading string
	pageContent string
	header      headers.Headers
}

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

func main() {
	server, err := server.Serve(port, handlerMain)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerMain(writer *response.Writer, req *request.Request) {
	log.Println("Handling request for", req.RequestLine.RequestTarget)
	page := pageData{}

	// Handle page
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		writer.WriteStatusLine(response.StatusCodeBadRequest)
		page.pageTitle = "400 Bad Request"
		page.pageHeading = "Bad Request"
		page.pageContent = "Your request honestly kinda sucked."
	case "/myproblem":
		writer.WriteStatusLine(response.StatusCodeInternalServerError)
		page.pageTitle = "500 Internal Server Error"
		page.pageHeading = "Internal Server Error"
		page.pageContent = "Okay, you know what? This one is on me."
	default:
		writer.WriteStatusLine(response.StatusCodeSuccess)
		page.pageTitle = "200 OK"
		page.pageHeading = "Success!"
		page.pageContent = "Your request was an absolute banger."
	}
	fullPage := writePage(page)
	page.header = response.GetDefaultHeaders(len(fullPage))
	page.header["Content-Type"] = "text/html"
	writer.WriteHeaders(page.header)
	writer.WriteBody([]byte(fullPage))
}

func writePage(page pageData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
  <html><head>
    <title>%s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>%s</p>
  </body>
</html>`, page.pageTitle, page.pageHeading, page.pageContent)
}
