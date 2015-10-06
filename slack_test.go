package main

import (
	"testing"
)

// test that the id return is Mario's Slack ID
func TestMarioId(t *testing.T) {
	type userData struct {
		Id       string
		Expected string
	}

	_, userId, err := startSlackRTM("xoxb-10454313842-D9I2egkjCpowGMORrHhr9k9d")

	if err != nil {
		t.Errorf("Error getting websocket and user id forlm Slack %v", err)
	}

	if userId != "U0ADC97QS" {
		t.Errorf("UserId error, expected U0ADC97QS, got %v", userId)
	}

}
