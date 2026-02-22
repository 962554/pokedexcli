package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	prompt         = "Pokedex > "
	unknownWarning = "Unknown command"
	argWarning     = "This command (%s) needs an arg.\n"
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
		command, ok := commands[word]
		if !ok {
			fmt.Println(unknownWarning)
			continue
		}
		var arg string
		argCommands := "explore catch inspect"
		if strings.Contains(argCommands, word) {
			if len(cleaned) < 2 {
				fmt.Printf(argWarning, word)
				continue
			}
			arg = cleaned[1]
		}
		err := command.callback(cfg, arg)
		if err != nil {
			fmt.Printf("Error running command: %s\n", word)
		}
	}
}
