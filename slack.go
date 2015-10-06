package main

import (
	"encoding/json"
	"fmt"
	//"golang.org/x/net/websocket"
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

// starts Slack real time messaging calling rtm.start
// return a websocket, a userID, an error
func startSlackRTM(token string) (string, string, error) {
	url := "https://slack.com/api/rtm.start?token=" + token

	res, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: Slack url is wrong")
		return "", "", err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error: could not read Slack response")
		return "", "", err
	}

	var connectionResponse slackResponse
	json.Unmarshal(body, &connectionResponse)

	if !connectionResponse.Ok {
		err = fmt.Errorf("Slack response was not ok:, %s", connectionResponse.Error)
		return "", "", err
	}

	successMsg := fmt.Sprintf("Success: slack returned a websocket [%s] and a userID [%s]", connectionResponse.Url, connectionResponse.Userdata.Id)

	return successMsg, connectionResponse.Userdata.Id, nil
}
