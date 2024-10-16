package library

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type Aof struct {
	File *os.File
	Rd   *bufio.Reader
	Mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		File: f,
		Rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.Mu.Lock()
			aof.File.Sync()
			aof.Mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.Mu.Lock()
	defer aof.Mu.Unlock()

	return aof.File.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.Mu.Lock()
	defer aof.Mu.Unlock()

	if len(value.Array) == 4 { // indicate there is TTL
		v := Value{}
		v.Typ = "bulk"
		v.Bulk = strconv.FormatInt(time.Now().Unix(), 10)
		value.Array = append(value.Array, v)
	}

	_, err := aof.File.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(fn func(value Value)) error {
	aof.Mu.Lock()
	defer aof.Mu.Unlock()

	aof.File.Seek(0, io.SeekStart)

	reader := NewResp(aof.File)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		fn(value)
	}

	return nil
}
