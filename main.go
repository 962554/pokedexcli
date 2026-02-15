package main

import (
	"bufio"
	"fmt"
	"os"
)

const prompt = "Pokedex > "

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]cliCommand{}

func commandExit() error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func usage() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func main() {
	// register exit command
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	// register exit command
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    usage,
	}

	scanner := bufio.NewScanner(os.Stdin)

	var userInput string

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
				_ = commands["exit"].callback()
			case "help":
				_ = commands["help"].callback()
			default:
				fmt.Println("Unknown command")
			}
		}
		if err := scanner.Err(); err != nil {
			_ = fmt.Errorf("%w", err)
		}
	}
}
