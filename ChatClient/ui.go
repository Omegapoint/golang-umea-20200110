package main

import (
	"bufio"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"os"
	"strings"
	"time"
)

var showPrompt = false

// userMessageRPLoop reads messages from the user and prints them to stdout as well as writing
// them to the provided channel.
func userMessageRPLoop(messages chan string) {
	fmt.Printf("\n\n=============================================\n")
	reader := bufio.NewReader(os.Stdin)
	var message string
	writePrompt()
	for {
		message, _ = reader.ReadString('\n')
		message = strings.TrimSpace(message)
		deletePreviousLine()
		printChatMessage(message, "me", true)
		messages <- message
	}
}

func deletePreviousLine() {
	fmt.Printf("\033[1A") // move cursor up n lines
	deleteCurrentLine()
}

func deleteCurrentLine() {
	fmt.Printf("\r\033[K") // delete to end of line
}

func writePrompt()  {
	showPrompt = true
	fmt.Print("Say something: ")
}

// printInfoMessage prints a given message as an 'info'-type message
func printInfoMessage(msg string) {
	printFunction(func() {
		fmt.Printf("\033[0;37m[%v] \033[1minfo:\033[0m %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	})
}

// printErrorMessage prints a given message as an 'error'-type message
func printErrorMessage(msg string) {
	printFunction(func() {
		fmt.Fprintf(os.Stderr, "\033[0;37m[%v] \033[0;31m\033[1merror:\033[0m\033[0m %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	})
}

// printChatMessage prints a given message `msg` as a chat message. The `user` parameter will be written
// as the users name, together with the message. Could i.e. be 'me' for the local client. If `fromLocal`
// true, the color of the name will be different from when the parameter is false.
func printChatMessage(msg string, user string, fromLocal bool) {
	printLocal := func() { fmt.Printf("\033[0;37m[%v]\033[0m \033[0;31m%s\033[0m: %s\n", time.Now().Format("2006-01-02 15:04:05"), user, msg) }
	printReceived := func() { fmt.Printf("\033[0;37m[%v]\033[0m \033[0;34m%s\033[0m: %s\n", time.Now().Format("2006-01-02 15:04:05"), user, msg) }
	if fromLocal {
		printFunction(printLocal)
	} else {
		printFunction(printReceived)
	}
}

func printFunction(f func()) {
	if showPrompt {
		deleteCurrentLine()
	}
	f()
	if showPrompt {
		writePrompt()
	}
}

func printJumboMessage(msg string) {
	myFigure := figure.NewFigure(msg, "", true)
	myFigure.Print()

	fmt.Println()
}
