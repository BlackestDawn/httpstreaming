package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	host = "127.0.0.1"
	port = "42069"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	lineReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		line, err := lineReader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
