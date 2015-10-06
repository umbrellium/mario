package main

import (
	"fmt"
	"log"
)

// TO DO: replace this global var with a ENVIRONEMT VARIABLE
const token string = "xoxb-10454313842-D9I2egkjCpowGMORrHhr9k9d"

func main() {
	res, _, err := startSlackRTM(token)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
