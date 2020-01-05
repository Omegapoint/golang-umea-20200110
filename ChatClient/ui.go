package main

import (
	"bufio"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"os"
	"strings"
	"time"
)

// userMessageRPLoop reads messages from the user and prints them to stdout as well as writing
// them to the provided channel.
func userMessageRPLoop(messages chan string) {
	fmt.Printf("\n\n=============================================\n")
	reader := bufio.NewReader(os.Stdin)
	var message string
	for {
		fmt.Print("Say something: ")
		message, _ = reader.ReadString('\n')
		fmt.Printf("\033[1A") // move cursor up one line
		fmt.Printf("\033[K") // delete to end of line
		printChatMessage(message, "me", true)
		messages <- strings.TrimSpace(message)
	}
}

func printInfoMessage(msg string) {
	fmt.Printf("\033[0;37m[%v] \033[1minfo:\033[0m %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func printErrorMessage(msg string) {
	fmt.Fprintf(os.Stderr, "\033[0;37m[%v] \033[0;31m\033[1merror:\033[0m\033[0m %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func printChatMessage(msg string, user string, fromLocal bool) {
	if fromLocal {
		fmt.Printf("\033[0;37m[%v]\033[0m \033[0;31m%s\033[0m: %s\n", time.Now().Format("2006-01-02 15:04:05"), user, msg)
	} else {
		fmt.Printf("\033[0;37m[%v]\033[0m \033[1;35m%s\033[0m: %s\n", time.Now().Format("2006-01-02 15:04:05"), user, msg)
	}
}

func printJumboMessage(msg string) {
	myFigure := figure.NewFigure(msg, "", true)
	myFigure.Print()

	fmt.Println()
}
