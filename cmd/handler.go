package main

import (
	"redisgow/cmd/library"
	"sync"
)

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var Handlers = map[string]func([]library.Value) library.Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

func ping(args []library.Value) library.Value {
	if len(args) == 0 {
		return library.Value{Typ: "string", Str: "PONG"}
	}

	return library.Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []library.Value) library.Value {
	if len(args) != 2 {
		return library.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	defer SETsMu.Unlock()

	return library.Value{Typ: "string", Str: "OK"}
}

func get(args []library.Value) library.Value {
	if len(args) != 1 {
		return library.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	defer SETsMu.RUnlock()

	if !ok {
		return library.Value{Typ: "null"}
	}

	return library.Value{Typ: "bulk", Bulk: value}
}
