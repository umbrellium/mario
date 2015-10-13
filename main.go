package main

import (
	"fmt"
	"github.com/umbrellium/mario/Godeps/_workspace/src/golang.org/x/net/websocket"
	"log"
	"os"
	"strings"
)

func main() {

	// instantiate slack
	var s Slack

	fmt.Println("Running Mario. Press ctrl+C to stop")

	// slack token must be set as environmet var or passed as command line
	token := os.Getenv("TOKEN")

	if token == "" {
		token = os.Args[1]
		if token == "" {
			log.Fatal("You must pass a token to connect to Slack")
		}
	}

	var ws *websocket.Conn
	// connect to slack
	ws, marioID, err := connectToSlack(token)

	// pass websocket to slack struct
	s.Socket = ws

	if err != nil {
		log.Fatal(err)
	}

	for {
		message, err := s.getMessage()

		if err != nil {
			log.Fatal(err)
		}

		// parse message and act accordingly
		if message.Type == "message" && strings.HasPrefix(message.Text, "<@"+marioID+">") {
			text := strings.TrimPrefix(message.Text, "<@"+marioID+"> ")
			text = strings.TrimSpace(text)

			messageHandled := false

			for _, task := range tasks {
				// we are using text to perform a reg ex and decide which method to call
				if task.Hear(&s, message, text) {
					messageHandled = true
					break
				}
			}

			// Mario cannot understand command
			if messageHandled == false {
				message.Text = `I don't understand what you are asking me to do.
Please ensure that your message doesn't contain any spelling mistake.
You can type '@mario help' to see a list of the available tasks I can perform.`
				err := s.postMessage(message)

				if err != nil {
					log.Fatal(err)
				}
			}

		}

	}
}
