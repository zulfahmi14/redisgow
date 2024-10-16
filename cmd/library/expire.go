package library

import (
	"math/rand"
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

func Sweep() {
	notExpiredKey := 0
	for key, value := range Ttls {
		if rand.Float64() < 0.5 { // 50% chance to check a key (adjust this probability as needed)
			if value <= time.Now().Unix() {
				SETsMu.Lock()
				delete(Ttls, key)
				delete(Data, key)
				SETsMu.Unlock()
			} else {
				notExpiredKey++
			}

			time.Sleep(10 * time.Millisecond)
		}

		if float64(notExpiredKey)/float64(len(Ttls)) > 0.75 {
			break
		}
	}
}
