package library

import (
	"sync"
	"time"
)

var Ttls = map[string]int64{}
var SETsMu = sync.RWMutex{}

func Expire(key string) {
	delete(Ttls, key)
}

func SetExpiry(key string, ttl int64) {
	Ttls[key] = time.Now().Unix() + ttl
}

func GetExpiry(key string) (int64, bool) {
	ttl, ok := Ttls[key]
	return ttl, ok
}
