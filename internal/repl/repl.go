// Package repl implements an interactive Read-Eval-Print Loop for the Pokedex CLI.
// It handles user input, command parsing, and execution.
package repl

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const (
	prompt         = "Pokedex > "
	unknownWarning = "Unknown command"
	argWarning     = "This command (%s) needs an arg.\n"
)

// CleanInput processes a raw string into a slice of lowercase words.
func CleanInput(text string) []string {
	lcText := strings.ToLower(text)

	return strings.Fields(lcText)
}

// RunRepl starts the interactive REPL loop, processing commands from standard input.
func RunRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		log.Print(prompt)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		err := scanner.Err()
		if err != nil {
			log.Println("reading input:", err)

			return
		}

		userInput := scanner.Text()

		cleaned := CleanInput(userInput)
		if len(cleaned) == 0 {
			continue
		}

		word := cleaned[0]
		commands := createCommands()

		command, exists := commands[word]
		if !exists {
			log.Println(unknownWarning)

			continue
		}

		var arg string

		argCommands := "explore catch inspect"
		if strings.Contains(argCommands, word) {
			if len(cleaned) < 2 {
				log.Printf(argWarning, word)

				continue
			}

			arg = cleaned[1]
		}

		err = command.callback(cfg, arg)
		if err != nil {
			log.Printf("Error running command: %s\n", word)
		}
	}
}
