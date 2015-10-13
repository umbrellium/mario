// A set of tasks that Mario can perform
// Each task must adhere to the Task interface
// Tasks must be added to the tasks slice in main.go

package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// The Task interface that Mario's commands must have
type Task interface {
	Hear(slack chatAgent, message Message, input string) bool
	Help(slack chatAgent, message Message) error
	getName() string
}

// register new task here
// by appending a task struct to the Task interface slice
var tasks []Task
var s2 Slack

func init() {
	tasks = append(tasks, Help{})
	tasks = append(tasks, Hello{})
	tasks = append(tasks, Say{})
}

// Hello Task
// Returns a simple Hello string
type Hello struct {
	Name string "hello"
}

// Hello Hear
// Returns true if the hello task is called
func (h Hello) Hear(slack chatAgent, message Message, input string) bool {
	// parse input and check if 'hello' is the first word
	r, err := regexp.Compile(`(?i)^\bhello\b`)
	if err != nil {
		fmt.Println("Error parsing Hello input")
	}

	if r.MatchString(input) {

		inputOptions := strings.Fields(input)

		if len(inputOptions) == 1 {
			Hello.say(h, slack, message)
			return true
		}

		if len(inputOptions) == 2 && inputOptions[1] == "help" {
			err := Hello.Help(h, slack, message)

			if err != nil {
				fmt.Println("Error posting message to slack")
				return false
			}
			return true
		}
		return false
	}
	return false
}

// Hello Help
// Returns a help string for the Hello struct
func (s Hello) Help(slack chatAgent, message Message) error {
	message.Text = `The <hello> command simply prints a hello message to Slack.
This command doesn't take any other options`
	err := slack.postMessage(message)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s Hello) getName() string {
	return "hello"
}

// Hello Say
// Posts a "Hello!" message to Slack
func (s Hello) say(slack chatAgent, message Message) error {
	message.Text = "Hello There!"
	err := slack.postMessage(message)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Help Task
// Returns help strings that explain how to use Mario's commands
type Help struct {
}

// Help Hear
// Returns true if the help task is called
func (s Help) Hear(slack chatAgent, message Message, input string) bool {
	r, err := regexp.Compile(`(?i)^\bhelp\b`)

	if err != nil {
		fmt.Println("Error parsing Help input")
	}

	if r.MatchString(input) {
		options := strings.Fields(input)

		if len(options) == 1 {
			// generice help
			err := s.Help(slack, message)

			if err != nil {
				fmt.Println("Error posting message to slack")
				return false
			}

			return true

		} else if len(options) == 2 && options[1] != "help" {
			// specific task help
			s.listCommands(slack, message, options[1])
			return true

		} else if len(options) == 2 && options[1] == "help" {
			// excetion: user typed "help" twice
			message.Text = `The <help> command doesn't take any argument.
Did you mean "@mario help" ?`
			err := slack.postMessage(message)
			if err != nil {
				log.Fatal(err)
			}
			return true
		}
		return false
	}
	return false
}

// Help Help
// post a generic help message to Slack
func (s Help) Help(slack chatAgent, message Message) error {

	message.Text = `Use this command to get an explanation about how to ask me 
to  perform a task.
Usage: 
- @mario help <command name>

Here is a list of the tasks I can currently perform:
`

	for _, t := range tasks {
		message.Text += "- " + t.getName() + fmt.Sprintf("\n")
	}

	err := slack.postMessage(message)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// getName Help
// return struct name
func (s Help) getName() string {
	return "help"
}

// listCommands lists the tasks that Mario can perform.
// Returns a message that will be posted to Slack
func (s Help) listCommands(slack chatAgent, message Message, command string) {
	commanHelp := command + " help"
	helpHandled := false

	for _, task := range tasks {
		if task.Hear(slack, message, commanHelp) {
			helpHandled = true
			break
		}
	}

	if helpHandled == false {
		message.Text = `I don't understand what you need help with.
Type "@mario help" for a list of tasks I can perfom.`
		err := slack.postMessage(message)

		if err != nil {
			log.Fatal(err)
		}
	}
}

// Say Task
// returns a custom string that will be posted to Slack
type Say struct {
	Name string "say"
}

// Hear Say
// Returns true if the say task is called
func (s Say) Hear(slack chatAgent, message Message, input string) bool {
	return false
}

// Help Say
// Returns a help string for the Say struct
func (s Say) Help(slack chatAgent, message Message) error {
	message.Text = `Use this command to tell Mario to send a message to Slack.
	Usage: 
	- @mario say "the message to post to Slack"`
	err := slack.postMessage(message)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// getName Say
// return struct name
func (s Say) getName() string {
	return "say"
}
