package main

import (
	"fmt"
	"net"
	"redisgow/cmd/library"
	"strings"
)

func main() {
	port := "6379"
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

	defer conn.Close() // close connection once finished

	for {
		resp := library.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := library.NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(library.Value{Typ: "string", Str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
