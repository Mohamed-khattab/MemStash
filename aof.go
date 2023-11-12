// contain persistent logs for every command executed in memory to be recompile when any fault happened
package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type aof struct {
	aofFile *os.File
	rd      *bufio.Reader
	mu      sync.RWMutex
}

// newAOF creates a new aof object and opens a file at the given path.
func newAOF(path string) (*aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	aof := &aof{
		aofFile: f,
		rd:      bufio.NewReader(f),
	}
	go func() {
		for {
			aof.mu.Lock()
			aof.aofFile.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()
	return aof, nil

}
// write writes the marshaled value to the aof file.
func (aof *aof) write(v Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
		// Write the marshaled value to the aof file.
	_, err := aof.aofFile.Write(v.Marshal())
	return err
	
}
// close closes the aof file and releases the associated resources.
func (aof *aof) close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.aofFile.Close()
}

