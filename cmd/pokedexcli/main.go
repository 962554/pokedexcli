// Package main provides a CLI tool for interacting with the PokeAPI.
// It allows users to browse location areas and manage a local Pokedex
// through an interactive REPL.
package main

import (
	"github.com/962554/pokedexcli/internal/repl"
)

func main() {
	repl.RunRepl()
}
