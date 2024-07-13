package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	port := "6000"
	fmt.Printf("Redisgow listening %s ", port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		buf := make([]byte, 1024)

		// read message from client
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		// giving response
		conn.Write([]byte("+OK\r\n"))
	}

	defer conn.Close() // close connection once finished
}