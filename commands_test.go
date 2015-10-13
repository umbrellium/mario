package main

import (
	"golang.org/x/net/websocket"
	"testing"
)

type FakeSlackChat struct{}

func (t *FakeSlackChat) getMessage(ws *websocket.Conn) (Message, error) {
	msg := new(Message)
	return *msg, nil
}

func (t *FakeSlackChat) postMessage(ws *websocket.Conn, msg Message) error {
	return nil
}

func TestHelloCommand(t *testing.T) {

	slack := new(FakeSlackChat)
	ws := new(websocket.Conn)
	msg := Message{1, "", "", ""}
	hello := new(Hello)

	type helloTestingStruct struct {
		input    string
		expected bool
	}

	helloTest := []helloTestingStruct{
		{"hello", true},
		{"hello ", true},
		{"Hello", true},
		{"", false},
		{"helloh", false},
		{"hello hello", false},
	}

	for _, tst := range helloTest {
		res := hello.Hear(slack, ws, msg, tst.input)
		if res != tst.expected {
			t.Errorf("Expected %q to return %q, got instead %q", tst.input, tst.expected, res)
		}
	}
}
