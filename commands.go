package main

import (
	"errors"
	"fmt"
	"math/rand"
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

var pokedex = map[string]pokeapi.Pokemon{}

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

func catchCommand(c *Config, pokemon string) error {
	_ = c
	const (
		commandMessage = "Throwing a Pokeball at %s...\n"
		caughtMessage  = "%s was caught!\n"
		escapedMessage = "%s escaped!\n"
	)
	resource, err := pokeapi.GetPokemon(pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("baseExperience: %d\n", resource.BaseExperience)

	count := 5
	difficulty := resource.BaseExperience / 10
	base := rand.Intn(difficulty)
	for range count {
		fmt.Printf(commandMessage, pokemon)
		try := rand.Intn(difficulty)
		if try != base {
			fmt.Printf(escapedMessage, pokemon)
		} else {
			break
		}
	}
	fmt.Printf(caughtMessage, pokemon)
	pokedex[pokemon] = resource
	return nil
}

func inspectCommand(c *Config, pokemon string) error {
	_ = c
	const (
		missingPokemon = "you have not caught that pokemon"
	)

	p, inPokedex := pokedex[pokemon]
	if !inPokedex {
		fmt.Println(missingPokemon)
		return errors.New(missingPokemon)
	}
	fmt.Println()
	fmt.Println("Name:", p.Name)
	fmt.Println("Height:", p.Height)
	fmt.Println("Weight:", p.Weight)
	fmt.Println("Stats:")
	for _, stat := range p.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typ := range p.Types {
		fmt.Printf("  - %s\n", typ.Type.Name)
	}
	return nil
}

func createCommands() map[string]cliCommand {
	commands := map[string]cliCommand{}

	// register commands
	commands["exit"] = newCliCommand("exit", "Exit the Pokedex", exitCommand)
	commands["help"] = newCliCommand("help", "Displays a help message", usageCommand)
	commands["map"] = newCliCommand("map", "Displays next 20 location-areas", mapCommand)
	commands["mapb"] = newCliCommand("mapb", "Displays previous 20 location-areas", mapbCommand)
	commands["explore"] = newCliCommand("explore", "Displays a list of all Pokemon within an area.", exploreCommand)
	commands["catch"] = newCliCommand("catch", "Catches a Pokemon and adds it to the Pokedex", catchCommand)
	commands["inspect"] = newCliCommand("inspect", "Prints the stats of a Pokemon in the Pokedex", inspectCommand)

	return commands
}
