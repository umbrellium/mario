// A set of tasks that Mario can perform
// Each task must adhere to the Task interface
// Tasks must be added to the tasks slice in main.go

package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"regexp"
	"strings"
)

type Task interface {
	Hear(ws *websocket.Conn, message Message, input string) bool
	Help(ws *websocket.Conn, message Message)
	getName() string
}

// register new task here
// by appending a task struct to the Task interface slice
var tasks []Task

func init() {
	tasks = append(tasks, Help{})
	tasks = append(tasks, Hello{})
	tasks = append(tasks, Say{})
}

//**** HELLO TASK ****//
// Returns a simple Hello string
type Hello struct {
	Name string "hello"
}

func (h Hello) Hear(ws *websocket.Conn, message Message, input string) bool {
	// parse input and check if 'hello' is the first word
	r, err := regexp.Compile(`(?i)^\bhello\b`)
	if err != nil {
		fmt.Println("Error parsing Hello input")
	}

	if r.MatchString(input) {

		inputOptions := strings.Fields(input)

		if len(inputOptions) == 1 {
			Hello.say(h, ws, message)
			return true
		} else {
			if inputOptions[1] == "help" {
				Hello.Help(h, ws, message)
				return true
			} else {
				return false
			}
		}
	} else {
		return false
	}
}

// Hello Help
// Returns a help string for the Hello struct
func (h Hello) Help(ws *websocket.Conn, message Message) {
	message.Text = `The <hello> command simply prints a hello message to Slack.
This command doesn't take any other options`
	err := postMessage(ws, message)

	if err != nil {
		log.Fatal(err)
	}
}

func (h Hello) getName() string {
	return "hello"
}

// Hello Say
// Posts a "Hello!" message to Slack
func (h Hello) say(ws *websocket.Conn, message Message) {
	message.Text = "Hello There!"
	err := postMessage(ws, message)

	if err != nil {
		log.Fatal(err)
	}
}

//**** HELP TASK ****//
// Returns help strings that explain how to use Mario's commands
type Help struct {
}

func (s Help) Hear(ws *websocket.Conn, message Message, input string) bool {
	r, err := regexp.Compile(`(?i)^\bhelp\b`)

	if err != nil {
		fmt.Println("Error parsing Help input")
	}

	if r.MatchString(input) {
		options := strings.Fields(input)

		if len(options) == 1 {
			// generice help
			s.Help(ws, message)
			return true
		} else if len(options) == 2 && options[1] != "help" {
			// specific task help
			s.listCommands(ws, message, options[1])
			return true
		} else if len(options) == 2 && options[1] == "help" {
			// excetion: user typed "help" twice
			message.Text = `The <help> commnad doesn't take any argument.
Did you mean "@mario help" ?`
			err := postMessage(ws, message)
			if err != nil {
				log.Fatal(err)
			}
			return true
		} else {
			// too many args
			return false
		}
	} else {
		return false
	}
}

func (h Help) Help(ws *websocket.Conn, message Message) {

	message.Text = `Use this command to get an explanation about how to ask me 
to  perform a task.
Usage: 
- @mario help <command name>

Here is a list of the tasks I can currently perform:
`

	for _, t := range tasks {
		message.Text += "- " + t.getName() + fmt.Sprintf("\n")
	}

	err := postMessage(ws, message)

	if err != nil {
		log.Fatal(err)
	}
}

func (h Help) getName() string {
	return "help"
}

// listCommands lists the tasks that Mario can perform.
// Returns a message that will be posted to Slack
func (h Help) listCommands(ws *websocket.Conn, message Message, command string) {
	commanHelp := command + " help"
	helpHandled := false

	for _, task := range tasks {
		if task.Hear(ws, message, commanHelp) {
			helpHandled = true
			break
		}
	}

	if helpHandled == false {
		message.Text = `I don't understand what you need help with.
Type "@mario help" for a list of tasks I can perfom.`
		err := postMessage(ws, message)

		if err != nil {
			log.Fatal(err)
		}
	}
}

//**** SAY TASK ****//
// returns a custom string that will be posted to Slack
type Say struct {
	Name string "say"
}

func (s Say) Hear(ws *websocket.Conn, message Message, input string) bool {
	return false
}

func (s Say) Help(ws *websocket.Conn, message Message) {
	message.Text = `Use this command to tell Mario to send a message to Slack.
	Usage: 
	- @mario say "the message to post to Slack"`
	err := postMessage(ws, message)

	if err != nil {
		log.Fatal(err)
	}
}

func (s Say) getName() string {
	return "say"
}
