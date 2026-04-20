package main

import (
	"bufio"
	"fmt"
)

const (
	ECHO  = 1
	P_REG = 2
	// Other message types
	R_ECHO  = 101
	R_P_REG = 102
)

type Message struct {
	ECHO  *string
	P_REG *string
	// Response
	R_ECHO  *string
	R_P_REG *byte // 0 Successes, 1 UnSuccesses
	// Other type here...
}

// Message format:
// -stream[0]: size
// ->stream[1:]: []byte
func readFromStream(streamRw *bufio.ReadWriter) ([]byte, error) {
	var err error
	// Read
	header, err := streamRw.ReadByte() // Block
	if err != nil {
		return nil, err
	}

	data, err := streamRw.Peek(int(header)) // Block
	if err != nil {
		return nil, err
	}

	_, err = streamRw.Discard(int(header))
	if err != nil {
		return nil, err
	}

	return data, err
}

func writeDataToStreamWithType(streamRw *bufio.ReadWriter, mType byte, data string) error {
	var err error
	// Write length
	err = streamRw.WriteByte(byte(len(data) + 1))
	if err != nil {
		return err
	}
	// Write type
	err = streamRw.WriteByte(mType)
	if err != nil {
		return err
	}
	// Write data
	_, err = streamRw.WriteString(data)
	if err != nil {
		return err
	}
	err = streamRw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func readMessageFromStream(streamRw *bufio.ReadWriter) (*Message, error) {
	data, err := readFromStream(streamRw)
	if err != nil {
		return nil, err
	}
	return parseMessage(data), nil
}

// [ 7  1  h e l l o o ]
func writeMessageToStream(streamRw *bufio.ReadWriter, message Message) error {
	if message.ECHO != nil {
		if err := writeDataToStreamWithType(streamRw, ECHO, *message.ECHO); err != nil {
			return err
		}
	} else if message.R_ECHO != nil {
		if err := writeDataToStreamWithType(streamRw, R_ECHO, *message.R_ECHO); err != nil {
			return err
		}
	}
	if message.P_REG != nil {
		if err := writeDataToStreamWithType(streamRw, P_REG, *message.P_REG); err != nil {
			return err
		}
	}
	if message.R_P_REG != nil {
		data := fmt.Sprintf("%d", *message.R_P_REG)
		if err := writeDataToStreamWithType(streamRw, R_P_REG, data); err != nil {
			return err
		}
	}
	return nil
}

func parseMessage(streamMessage []byte) *Message {
	switch streamMessage[0] {
	case ECHO:
		var st = string(streamMessage[1:])
		return &Message{ECHO: &st}
	case R_ECHO:
		var st = string(streamMessage[1:])
		return &Message{R_ECHO: &st}
	case P_REG:
		var st = string(streamMessage[1:])
		return &Message{P_REG: &st}
	case R_P_REG:
		var st = streamMessage[1]
		return &Message{R_P_REG: &st}
	default:
		return nil
	}
}
