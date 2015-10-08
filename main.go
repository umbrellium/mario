package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// slack token must be set as environmet var or passed as command line
	token := os.Getenv("TOKEN")
	if token == "" {
		token = os.Args[1]
		if token == "" {
			log.Fatal("You must pass a token to connect to Slack")
		}
	}

	// connect to slack
	websocket, marioId, err := connectToSlack(token)

	if err != nil {
		log.Fatal(err)
	}

	// start loop
	for {

		// get messages
		message, err := getMessage(websocket)

		if err != nil {
			log.Fatal(err)
		}

		// parse message and act accordingly
		if message.Type == "message" {
			// check if Mario was metioned
			msg_slice := strings.Fields(message.Text)
			if msg_slice[0] == "<@"+marioId+">" {

				//get command sent to mario
				command := msg_slice[1]

				switch command {
				case "help":
					go func(m Message) {
						message.Text = "Hello!"
						err := postMessage(websocket, message)

						if err != nil {
							log.Fatal(err)
						}

						fmt.Println("Posting message")
					}(message)
				default:
					fmt.Println("nothing was passed")
				}
			}
		}

	}
}
