package main

import (
	"encoding/json"
	"fmt"
	"github.com/umbrellium/mario/Godeps/_workspace/src/golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
)

type chatAgent interface {
	// getMessage(ws *websocket.Conn) (Message, error)
	// postMessage(ws *websocket.Conn, msg Message) error
	getMessage() (Message, error)
	postMessage(msg Message) error
}

type Slack struct {
	Socket *websocket.Conn
}

type slackResponse struct {
	Ok       bool   `json:"ok"`
	Error    string `json:"error"`
	Url      string `json:"url"`
	Userdata userID `json:"self"`
}

type userID struct {
	Id string `json:"id"`
}

// Message struct use to generate Slack messages
type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

var counter uint64

// ConnectToSlack starts Slack real time messaging and opens a websocket
// Returns a websocket, a userID, an error
func connectToSlack(token string) (websocket.Conn, string, error) {
	url := "https://slack.com/api/rtm.start?token=" + token

	// connect to rtm
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// store get response
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	// assign get response to slackResponse struct
	var connectionResponse slackResponse
	json.Unmarshal(body, &connectionResponse)

	if !connectionResponse.Ok {
		log.Fatal(err)
	}

	// connect to slack
	socket, err := websocket.Dial(connectionResponse.Url, "", "https://api.slack.com/")

	if err != nil {
		err = fmt.Errorf("Error: cannot open slack websocket")
		return socket, "", err
	}

	return socket, connectionResponse.Userdata.Id, nil
}

// GetMessage listens to Slack messages
// Returns the message or an error
func (s *Slack) getMessage() (Message, error) {

	// the message to return
	var msg Message
	err := websocket.JSON.Receive(s.socket, &msg)

	if err != nil {
		fmt.Errorf("Error: cannot get message")
		return msg, err
	}

	return msg, err
}

// PostMessage publishes a message on Slack
// Returns an error if it couldn't complete the operation
func (s *Slack) postMessage(msg Message) error {
	msg.Id = atomic.AddUint64(&counter, 1)
	err := websocket.JSON.Send(s.socket, msg)
	return err
}
