.PHONY: client
client:
	caddy run

.PHONY: server
server:
	go run main.go midi.go stream.go
