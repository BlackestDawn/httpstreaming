package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Port     int
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

func Serve(port int, handler Handler) (*Server, error) {
	// Initialize server
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Port:     port,
		listener: listener,
		handler:  handler,
	}

	go func() {
		s.listen()
	}()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle connection
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	log.Println("New connection from", conn.RemoteAddr())
	// Read request
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println("Error reading request:", err)
		WriteError(conn, &HandlerError{
			Code:    response.StatusCodeInternalServerError,
			Message: err.Error(),
		})
		return
	}

	// Handle request
	writer := response.NewWriter(conn)
	s.handler(writer, req)
}

func WriteError(w io.Writer, err *HandlerError) {
	response.WriteStatusLine(w, err.Code)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(err.Message)))
	w.Write([]byte(err.Message))
}
