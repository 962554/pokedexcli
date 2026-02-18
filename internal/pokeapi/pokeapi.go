package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Next     string
	Previous string
}

type NamedAPIResource struct {
	Count int
	Config
	Results []struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	}
}

// GetLocationAreas returns all location-areas from PokeAPI
func GetLocationAreas(url string) (NamedAPIResource, error) {
	res, err := http.Get(url)
	if err != nil {
		return NamedAPIResource{}, fmt.Errorf("http.Get failed: %w", err)
	}
	defer res.Body.Close()

	var resource NamedAPIResource
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&resource); err != nil {
		return NamedAPIResource{}, fmt.Errorf("decoder.Decode failed: %w", err)
	}

	return resource, nil
}

func ExtractConfig(resource NamedAPIResource) *Config {
	return &resource.Config
}
