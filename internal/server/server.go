package server

import (
	"bytes"
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
	running  atomic.Bool
	listener net.Listener
	handler  Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

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
		running:  atomic.Bool{},
		listener: listener,
		handler:  handler,
	}
	s.running.Store(true)

	go func() {
		s.listen()
	}()

	return s, nil
}

func (s *Server) Close() error {
	s.running.Store(false)
	return s.listener.Close()
}

func (s *Server) listen() error {
	for conn, err := s.listener.Accept(); err == nil; conn, err = s.listener.Accept() {
		if !s.running.Load() {
			break
		}

		// Handle connection
		go func() {
			s.handle(conn)

		}()
	}
	return nil
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

	log.Println("Request received:", req.RequestLine.Method, req.RequestLine.RequestTarget)
	// Handle request
	buf := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buf, req)
	if handlerErr != nil {
		log.Println("Error handling request:", handlerErr.Message)
		WriteError(conn, handlerErr)
		return
	}

	log.Println("Request handled successfully")
	// Send response
	head := response.GetDefaultHeaders(buf.Len())
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	response.WriteHeaders(conn, head)
	conn.Write(buf.Bytes())
}

func (s *Server) IsRunning() bool {
	return s.running.Load()
}

func WriteError(w io.Writer, err *HandlerError) {
	response.WriteStatusLine(w, err.Code)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(err.Message)))
	w.Write([]byte(err.Message))
}
