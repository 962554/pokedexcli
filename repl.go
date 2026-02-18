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
)

func cleanInput(text string) []string {
	lcText := strings.ToLower(text)
	return strings.Fields(lcText)
}

func runRepl() {

	scanner := bufio.NewScanner(os.Stdin)

	var userInput string
	commands := createCommands()
	control = Config{
		Next:     mapEndpoint,
		Previous: "",
	}

	for {
		fmt.Print(prompt)
		if scanner.Scan() {
			userInput = scanner.Text()
			cleaned := cleanInput(userInput)
			if len(cleaned) == 0 {
				continue
			}
			command := cleaned[0]
			switch command {
			case "exit":
				_ = commands["exit"].callback(&control)
			case "help":
				_ = commands["help"].callback(&control)
			case "map":
				_ = commands["map"].callback(&control)
			case "mapb":
				_ = commands["mapb"].callback(&control)
			default:
				fmt.Println("Unknown command")
			}
		}
		if err := scanner.Err(); err != nil {
			_ = fmt.Errorf("%w", err)
		}
	}
}
