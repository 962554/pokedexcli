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

func newCliCommand(name, desc string, f func(*Config, string) error) cliCommand {
	return cliCommand{
		name:        name,
		description: desc,
		callback:    f,
	}
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

	}

	}


	// register commands
	commands["exit"] = newCliCommand("exit", "Exit the Pokedex", exitCommand)
	commands["help"] = newCliCommand("help", "Displays a help message", usageCommand)
	commands["map"] = newCliCommand("map", "Displays next 20 location-areas", mapCommand)
	commands["mapb"] = newCliCommand("mapb", "Displays previous 20 location-areas", mapbCommand)
	commands["explore"] = newCliCommand("explore", "Displays a list of all Pokemon within an area.", exploreCommand)
	commands["catch"] = newCliCommand("catch", "Catches a Pokemon and adds it to the Pokedex", catchCommand)

	return commands
}
