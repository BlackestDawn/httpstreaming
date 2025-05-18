package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

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

func handlerMain(w io.Writer, req *request.Request) *server.HandlerError {
	log.Println("Handling request for", req.RequestLine.RequestTarget)
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			Code:    response.StatusCodeBadRequest,
			Message: "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			Code:    response.StatusCodeInternalServerError,
			Message: "Woopsie, my bad\n",
		}
	default:
		fmt.Fprintf(w, "All good, frfr\n")
		return nil
	}
}
