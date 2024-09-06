package main

import (
	"bytes"
	"io"
	"sync"
	"time"
)

const (
	BUFFERSIZE = 8192
	DELAY      = 1000 // this needs to be at least as quick as the bpm
)

type Connection struct {
	bufferChannel chan []byte
	buffer        []byte
}

type ConnectionPool struct {
	connectionMap map[*Connection]struct{}
	mu            sync.Mutex
}

func (cp *ConnectionPool) AddConnection(connection *Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.connectionMap[connection] = struct{}{}
}

func (cp *ConnectionPool) DeleteConnection(connection *Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	delete(cp.connectionMap, connection)
}

func (cp *ConnectionPool) Broadcast(buffer []byte) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for connection := range cp.connectionMap {
		copy(connection.buffer, buffer)

		select {
		case connection.bufferChannel <- connection.buffer:
		default:
		}
	}
}

func NewConnectionPool() *ConnectionPool {
	connectionMap := make(map[*Connection]struct{})

	return &ConnectionPool{connectionMap: connectionMap}
}

func stream(connectionPool *ConnectionPool, content []byte) {
	buffer := make([]byte, BUFFERSIZE)

	for {
		clear(buffer)
		tempfile := bytes.NewReader(content)
		ticker := time.NewTicker(time.Millisecond * DELAY)

		for range ticker.C {
			_, err := tempfile.Read(buffer)

			if err == io.EOF {
				ticker.Stop()
				break
			}

			connectionPool.Broadcast(buffer)
		}
	}
}
