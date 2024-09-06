package main

import (
	"log"
	"net/http"
)

func main() {
	connPool := NewConnectionPool()
	content := makeMIDI()

	go stream(connPool, content)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "audio/midi")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("Could not create flusher")
		}

		connection := &Connection{
			bufferChannel: make(chan []byte),
			buffer:        make([]byte, BUFFERSIZE),
		}
		connPool.AddConnection(connection)
		log.Printf("%s has connected to the audio stream\n", r.Host)

		for {
			buff := <-connection.bufferChannel
			_, err := w.Write(buff)

			if err != nil {
				connPool.DeleteConnection(connection)
				log.Printf("%s's connection to the audio stream has been closed\n", r.Host)

				return
			}

			flusher.Flush()
			clear(connection.buffer)
		}
	})

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
