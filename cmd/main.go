package main

import (
	"fmt"
	"net"
	"redisgow/cmd/library"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := "6379"
	fmt.Printf("Redisgow listening %s ", port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := library.NewAof("aof/database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value library.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		if len(args) == 4 { // if TTL has been passed, do not insert!
			v := args[len(args)-1]   // timestamp
			ttl := args[len(args)-2] // ttl
			timestamp, _ := strconv.ParseInt(v.Bulk, 10, 64)
			ttlNum, _ := strconv.ParseInt(ttl.Bulk, 10, 64)
			if v.Typ == "bulk" && timestamp+ttlNum <= time.Now().Unix() {
				return
			}

			ttl.Bulk = strconv.FormatInt(timestamp+ttlNum-time.Now().Unix(), 10) // update TTL based on current date
			args[len(args)-2] = ttl
			args = args[:len(args)-1] // remove the timestamp info
		}

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				library.Sweep()
			}
		}
	}()

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

		if command == "SET" {
			aof.Write(value)
		}

		result := handler(args)

		writer.Write(result)
	}
}
