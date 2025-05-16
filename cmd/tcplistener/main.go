package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const portNum = "42069"

func main() {
	socket, err := net.Listen("tcp", ":"+portNum)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connection opened")
		defer conn.Close()

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(req)
		fmt.Println("connection closed")
	}
}
