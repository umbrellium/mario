package main

import (
	"testing"
)

var slack FakeSlackChat
var msg = Message{1, "", "", ""}

// FakeSlackChat is a fake slack requests struct
// it implements the chatAgent interface defined in slack.go
// so that it can replace the real slack requests when testing
type FakeSlackChat struct{}

// fakeSlackChat implement get message
func (t *FakeSlackChat) getMessage() (Message, error) {
	msg := new(Message)
	return *msg, nil
}

// fakeSlackChat implement post message
func (t *FakeSlackChat) postMessage(msg Message) error {
	return nil
}

// TestHelloCommand tests responses from the <hello> command
func TestHelloHearCommand(t *testing.T) {
	hello := new(Hello)

	type helloTestingStruct struct {
		input    string
		expected bool
	}

	helloHearTest := []helloTestingStruct{
		{"hello", true},
		{"hello ", true},
		{"Hello", true},
		{"", false},
		{"helloh", false},
		{"hello hello", false},
	}

	for _, tst := range helloHearTest {
		res := hello.Hear(&slack, msg, tst.input)
		if res != tst.expected {
			t.Errorf("Expected %q to return %q, got %q instead", tst.input, tst.expected, res)
		}
	}
}

func TestHelloHearHelp(t *testing.T) {
	hello := new(Hello)

	type helloTestingStruct struct {
		input    string
		expected bool
	}

	helloHearHelp := []helloTestingStruct{
		{"hello help", true},
		{"hello help help", false},
		{"hello help Help", false},
		{"hello help hello", false},
		{"hello helpme", false},
		{"hello Help", false},
	}

	for _, tst := range helloHearHelp {
		res := hello.Hear(&slack, msg, tst.input)
		if res != tst.expected {
			t.Errorf("Expected %q to return %q, got %q instead", tst.input, tst.expected, res)
		}
	}
}

func TestHelloSay(t *testing.T) {
	hello := new(Hello)

	err := hello.say(&slack, msg)

	if err != nil {
		t.Errorf("Expected Hello.say to return no error")
	}
}
