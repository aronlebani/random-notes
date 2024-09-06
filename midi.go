package main

import (
	"gitlab.com/gomidi/midi/v2"
)

func appendMessage(bytes []byte, message midi.Message) []byte {
	for _, b := range message.Bytes() {
		bytes = append(bytes, b)
	}

	return bytes
}

func makeMIDI() []byte {
	var bytes []byte

	// Octaves seem to be zero-indexed, so middle C is C(5), not C(4)
	bytes = appendMessage(bytes, midi.NoteOn(0, midi.C(5), 120))
	bytes = appendMessage(bytes, midi.NoteOn(0, midi.Ab(3), 120))
	bytes = appendMessage(bytes, midi.NoteOff(0, midi.C(5)))
	bytes = appendMessage(bytes, midi.NoteOff(0, midi.Ab(3)))

	return bytes
}
