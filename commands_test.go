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

	test1 := hello.Hear(slack, ws, msg, "hello")

	if test1 != true {
		t.Fatalf("Expected to hear 'hello', returned %q", test1)
	}
}
