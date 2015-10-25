// A set of tasks that Mario can perform
// Each task must adhere to the Task interface
// Tasks must be added to the tasks slice in main.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	tasks = append(tasks, Wercker{})
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
	message.Text = "Yo"
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
}

// Hear Say
// Returns true if the say task is called
func (s Say) Hear(slack chatAgent, message Message, input string) bool {
	r, err := regexp.Compile(`(?i)^\bsay\b`)

	if err != nil {
		fmt.Println("Error parsing Help input")
	}

	if r.MatchString(input) {

		options := strings.Fields(input)

		if len(options) == 1 || options[1] == "help" {
			err := Say.Help(s, slack, message)

			if err != nil {
				fmt.Println("Error parsing Hear option")
				return false
			}
			return true
		}
	}
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

// Wercker struct
// performs Wercker realted tasks (e.g. list apps, deploy app etc)
type Wercker struct {
}

type werkerApps struct {
	Name string `json:name`
}

func (s Wercker) Hear(slack chatAgent, message Message, input string) bool {
	patter, err := regexp.Compile(`^\blist apps\b`)

	if err != nil {
		fmt.Println("Error parsing Help input")
	}

	if patter.MatchString(input) {
		options := strings.Fields(input)

		if len(options) == 2 {
			res, err := Wercker.connectToAPI(s, "applications")
			if err != nil {
				fmt.Println("Error connecting to the Wercker API")
				return false
			}

			Wercker.listApps(s, res, slack, message)
			return true
		}

		if options[2] == "help" {
			// call help
			err := Wercker.Help(s, slack, message)
			if err != nil {
				fmt.Println("Error calling Wercker Help")
				return false
			}
			return true
		}
	}
	return false
}

func (s Wercker) connectToAPI(endpoint string) (*http.Response, error) {
	wtoken := os.Getenv("WERCKER_TOKEN")
	if wtoken == "" {
		wtoken = os.Args[2]
		// NOTE: token can be an empty string
		// Wercker will retrun only public apps
	}

	var url string

	switch endpoint {
	case "applications":
		url = "https://app.wercker.com/api/v3/applications/umbrellium?token=" + wtoken
	case "builds":
		url = "https://app.wercker.com/api/v3/builds/"
	case "deploy":
		url = "https://app.wercker.com/api/v3/deploys/"
	}

	res, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: problem talking to Wercker API")
		return nil, err
	}

	return res, nil
}

// Wercker listApps
// prints a list of Umbrellium apps that are currently available on Wercker
func (s Wercker) listApps(httpRes *http.Response, slack chatAgent, message Message) error {

	var availbaleApps []werkerApps

	// parse response
	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		fmt.Println("Error: problem reading response from Wercker API")
		return err
	}

	json.Unmarshal(body, &availbaleApps)

	message.Text = "The following apps are currently available on Wercker: \n"

	// print response to slack
	for _, app := range availbaleApps {
		fmt.Println(app)
		message.Text += app.Name + fmt.Sprintf("\n")
	}

	err2 := slack.postMessage(message)
	if err2 != nil {
		fmt.Println("Error: problem posting message to Slack")
		return err2
	}

	return nil
}

func (s Wercker) Help(slack chatAgent, message Message) error {
	message.Text = `<list apps> will list the Umbrellium applications currently available on Wercker. 
This command does not take any option. 
`

	err := slack.postMessage(message)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (s Wercker) getName() string {
	return "list apps"
}
