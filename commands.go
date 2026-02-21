package main

import (
	"fmt"
	"os"

	"github.com/962554/pokedexcli/internal/pokeapi"
)

// Config maintains the pagination state for API requests, storing the
// URLs for the next and previous sets of results.
type Config struct {
	next     *string
	previous *string
}

var cfg = &Config{
	next: &pokeapi.MapEndpoint,
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config) error
}

func (c cliCommand) String() string {
	return fmt.Sprintf("name: %s, description: %s", c.name, c.description)
}

func exitCommand(c *Config) error {
	_ = c
	defer os.Exit(0)
	fmt.Println(exitMessage)
	return nil
}

func usageCommand(c *Config) error {
	_ = c
	fmt.Println(welcomeMessage)
	fmt.Println(usageMessage)
	fmt.Println()

	for k, v := range createCommands() {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func mapCommand(c *Config) error {
	if c.next == nil {
		fmt.Println("you're on the last page")
		return nil
	}

	resource, err := pokeapi.GetLocationAreas(*c.next)
	if err != nil {
		return err
	}
	c.next = resource.Next
	c.previous = resource.Previous

	for _, result := range resource.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func mapbCommand(c *Config) error {
	if c.previous == nil {
		fmt.Println("you're on the first page")
		c.next = &pokeapi.MapEndpoint
		return nil
	}
	resource, err := pokeapi.GetLocationAreas(*c.previous)
	if err != nil {
		return err
	}
	c.next = resource.Next
	c.previous = resource.Previous

	for _, result := range resource.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func createCommands() map[string]cliCommand {

	var commands = map[string]cliCommand{}

	// register exit command
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    exitCommand,
	}

	// register exit command
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    usageCommand,
	}

	// register map command
	commands["map"] = cliCommand{
		name:        "map",
		description: "Displays next 20 location-areas",
		callback:    mapCommand,
	}

	// register mapb command
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays previous 20 location-areas",
		callback:    mapbCommand,
	}
	return commands
}
