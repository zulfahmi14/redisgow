package main

import (
	"redisgow/cmd/library"
	"strconv"
	"sync"
	"time"
)

var SETsMu = sync.RWMutex{}

var Handlers = map[string]func([]library.Value) library.Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"TTL":  ttl,
}

func ping(args []library.Value) library.Value {
	if len(args) == 0 {
		return library.Value{Typ: "string", Str: "PONG"}
	}

	return library.Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []library.Value) library.Value {
	if len(args) < 2 || len(args) > 3 {
		return library.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	defer SETsMu.Unlock()

	if len(args) == 3 {
		ttl, err := strconv.ParseInt(args[2].Bulk, 10, 64)
		if err == nil {
			library.SetExpiry(key, ttl)
		}
	}
	library.SetData(key, value)

	return library.Value{Typ: "string", Str: "OK"}
}

func get(args []library.Value) library.Value {
	if len(args) != 1 {
		return library.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	defer SETsMu.RUnlock()

	value, ok := library.GetData(key)
	if !ok {
		return library.Value{Typ: "null"}
	}

	ttl, ttlOk := library.GetExpiry(key)
	if ttlOk {
		if ttl-time.Now().Unix() <= 0 {
			library.Expire(key)
			library.DeleteData(key)
			return library.Value{Typ: "null"}
		}
	}

	return library.Value{Typ: "bulk", Bulk: value}
}

func ttl(args []library.Value) library.Value {
	if len(args) != 1 {
		return library.Value{Typ: "error", Str: "ERR wrong number of arguments for 'ttl' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	defer SETsMu.RUnlock()

	value, ok := library.GetExpiry(key)

	if !ok {
		return library.Value{Typ: "null"}
	}

	if value-time.Now().Unix() <= 0 {
		library.Expire(key)
		library.DeleteData(key)
		return library.Value{Typ: "int", Num: -1}
	}

	return library.Value{Typ: "int", Num: value - time.Now().Unix()}
}
