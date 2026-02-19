package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	prompt         = "Pokedex > "
	welcomeMessage = "Welcome to Pokedex!"
	usageMessage   = "Usage:"
	exitMessage    = "Closing the Pokedex... Goodbye!"
	unknownMessage = "Unknown command"
)

func cleanInput(text string) []string {
	lcText := strings.ToLower(text)
	return strings.Fields(lcText)
}

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		userInput := scanner.Text()
		cleaned := cleanInput(userInput)
		if len(cleaned) == 0 {
			continue
		}
		word := cleaned[0]
		commands := createCommands()
		if command, ok := commands[word]; !ok {
			fmt.Println(unknownMessage)
			continue
		} else {
			err := command.callback(cfg)
			if err != nil {
				fmt.Printf("Error running command: %s", word)
			}
		}
	}
}
