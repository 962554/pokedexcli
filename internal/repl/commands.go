package repl

import (
	"errors"
	"fmt"
	"log"
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

func (cmd cliCommand) String() string {
	return fmt.Sprintf("name: %s, description: %s", cmd.name, cmd.description)
}

func exitCommand(cfg *Config, arg string) error {
	_, _ = cfg, arg

	defer os.Exit(0)
	log.Println(exitMessage)

	return nil
}

func usageCommand(cfg *Config, arg string) error {
	_, _ = cfg, arg
	log.Println(welcomeMessage)
	log.Println(usageMessage)
	log.Println()

	for k, v := range createCommands() {
		log.Printf("%s: %s\n", k, v.description)
	}

	return nil
}

func mapCommand(cfg *Config, arg string) error {
	_ = arg
	if cfg.next == nil {
		log.Println("you're on the last page")

		return nil
	}

	resource, err := pokeapi.GetLocationAreas(*cfg.next)
	if err != nil {
		return fmt.Errorf("mapCommand: GetLocationArea failed: %w", err)
	}
	cfg.next = resource.Next
	cfg.previous = resource.Previous

	for _, result := range resource.Results {
		log.Println(result.Name)
	}

	return nil
}

func mapbCommand(cfg *Config, arg string) error {
	_ = arg
	if cfg.previous == nil {
		log.Println("you're on the first page")
		cfg.next = &pokeapi.MapEndpoint

		return nil
	}
	resource, err := pokeapi.GetLocationAreas(*cfg.previous)
	if err != nil {
		return fmt.Errorf("mapbCommand: GetLocationAreas failed: %v", err)
	}
	cfg.next = resource.Next
	cfg.previous = resource.Previous

	for _, result := range resource.Results {
		log.Println(result.Name)
	}

	return nil
}

func exploreCommand(cfg *Config, area string) error {
	_ = cfg
	resource, err := pokeapi.GetLocation(area)
	if err != nil {
		return fmt.Errorf("exploreCommand: GetLocation failed: %v", err)
	}
	log.Printf("Exploring %s...\n", area)
	log.Println("Found Pokemon: ")
	for _, result := range resource.PokemonEncounters {
		log.Printf(" - %s\n", result.Pokemon.Name)
	}

	return nil
}

func catchCommand(cfg *Config, pokemon string) error {
	_ = cfg
	const (
		commandMessage = "Throwing a Pokeball at %s...\n"
		caughtMessage  = "%s was caught!\n"
		escapedMessage = "%s escaped!\n"
	)
	resource, err := pokeapi.GetPokemon(pokemon)
	if err != nil {
		return fmt.Errorf("catchCommand: GetPokemon failed: %v", err)
	}

	log.Printf("baseExperience: %d\n", resource.BaseExperience)

	count := 5
	difficulty := resource.BaseExperience / 10
	base := rand.Intn(difficulty)
	for range count {
		log.Printf(commandMessage, pokemon)
		try := rand.Intn(difficulty)
		if try != base {
			log.Printf(escapedMessage, pokemon)
		} else {
			break
		}
	}
	log.Printf(caughtMessage, pokemon)
	pokedex[pokemon] = resource

	return nil
}

func inspectCommand(cfg *Config, pokemon string) error {
	_ = cfg
	const (
		missingPokemon = "you have not caught that pokemon"
	)

	pokemo, inPokedex := pokedex[pokemon]
	if !inPokedex {
		log.Println(missingPokemon)
		return errors.New(missingPokemon)
	}
	log.Println()
	log.Println("Name:", pokemo.Name)
	log.Println("Height:", pokemo.Height)
	log.Println("Weight:", pokemo.Weight)
	log.Println("Stats:")
	for _, stat := range pokemo.Stats {
		log.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	log.Println("Types:")
	for _, typ := range pokemo.Types {
		log.Printf("  - %s\n", typ.Type.Name)
	}

	return nil
}

func pokedexCommand(cfg *Config, arg string) error {
	_, _ = cfg, arg
	const (
		noPokemon      = "you don't have an Pokemon"
		pokedexHeading = "Your Pokedex:"
	)
	if len(pokedex) == 0 {
		log.Println(noPokemon)
		return errors.New(noPokemon)
	}
	log.Println()
	log.Println(pokedexHeading)
	for pokemon := range pokedex {
		log.Printf(" - %s\n", pokemon)
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
	commands["pokedex"] = newCliCommand("pokedex", "Prints the names of all Pokemon in the Pokedex", pokedexCommand)

	return commands
}
