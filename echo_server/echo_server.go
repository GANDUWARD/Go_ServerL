package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func worker(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 512)
	for {
		size, err := conn.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received %v bytes from %v\n", size, conn.RemoteAddr())
		size, err = conn.Write(b[0:size])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Written %v bytes to %v\n", size, conn.RemoteAddr())
	}
}

func main() {
	listner, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {

		log.Fatal(err)
	}
	fmt.Printf("Listering on %v\n", listner.Addr())
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Accepted connection to %v from %v\n", conn.LocalAddr(), conn.RemoteAddr())
		go worker(conn)
	}
}
