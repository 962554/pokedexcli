package main

import (
	"fmt"
	"os"

	"github.com/962554/pokedexcli/internal/pokeapi"
)

const (
	welcomeMessage = "Welcome to Pokedex!"
	usageMessage   = "Usage:"
	exitMessage    = "Closing the Pokedex... Goodbye!"
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
	callback    func(*Config, string) error
}

func (c cliCommand) String() string {
	return fmt.Sprintf("name: %s, description: %s", c.name, c.description)
}

func exitCommand(c *Config, area string) error {
	_, _ = c, area

	defer os.Exit(0)
	fmt.Println(exitMessage)
	return nil
}

func usageCommand(c *Config, area string) error {
	_, _ = c, area
	fmt.Println(welcomeMessage)
	fmt.Println(usageMessage)
	fmt.Println()

	for k, v := range createCommands() {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func mapCommand(c *Config, area string) error {
	_ = area
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

func mapbCommand(c *Config, area string) error {
	_ = area
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

func exploreCommand(c *Config, area string) error {
	_ = c
	resource, err := pokeapi.GetLocation(area)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon: ")
	for _, result := range resource.PokemonEncounters {
		fmt.Printf(" - %s\n", result.Pokemon.Name)
	}
	return nil
}

func createCommands() map[string]cliCommand {
	commands := map[string]cliCommand{}

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

	// register mapb command
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Displays a list of all Pokemon within an area.",
		callback:    exploreCommand,
	}
	return commands
}
