package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
)

type slackResponse struct {
	Ok       bool   `json:"ok"`
	Error    string `json:"error"`
	Url      string `json:"url"`
	Userdata userID `json:"self"`
}

type userID struct {
	Id string `json:"id"`
}

type message struct {
	Id      int    `json:id`
	Type    string `json:type`
	Channel string `json:channel`
	Text    string `json:text`
}

// ConnectToSlack starts Slack real time messaging and opens a websocket
// Returns a websocket, a userID, an error
func connectToSlack(token string) (string, string, error) {
	url := "https://slack.com/api/rtm.start?token=" + token

	// connect to rtm
	res, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: Slack url is wrong")
		return "", "", err
	}

	// store get response
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error: cannot read Slack response")
		return "", "", err
	}

	// assign get response to slackResponse struct
	var connectionResponse slackResponse
	json.Unmarshal(body, &connectionResponse)

	if !connectionResponse.Ok {
		err = fmt.Errorf("Slack response was not ok:, %s", connectionResponse.Error)
		return "", "", err
	}

	// connect to slack
	socket, err := websocket.Dial(connectionResponse.Url, "", "https://api.slack.com/")

	if err != nil {
		err = fmt.Errorf("Error: cannot open slack websocket")
		return "", "", err
	}

	return socket, connectionResponse.Userdata.Id

	successMsg := fmt.Sprintf("Success: slack returned a websocket [%s] and a userID [%s]", connectionResponse.Url, connectionResponse.Userdata.Id)
	return successMsg, connectionResponse.Userdata.Id, nil
}

// getMessage listens to Slack messages
// returns ...
func getMessage(webSocket string, msg string) string {

}

// postMessage publishes a message on Slack
// returns ...
func postMessage(webSocket string) string {

}
