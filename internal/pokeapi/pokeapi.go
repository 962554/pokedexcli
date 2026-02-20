package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	MapEndpoint string = "https://pokeapi.co/api/v2/location-area/"
)

type NamedAPIResource struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// GetLocationAreas returns all location-areas from PokeAPI
func GetLocationAreas(url string) (NamedAPIResource, error) {
	res, err := http.Get(url)
	if err != nil {
		return NamedAPIResource{}, fmt.Errorf("http.Get failed: %w", err)
	}
	defer res.Body.Close()

	var resource NamedAPIResource

	if err := json.Unmarshal(data, &resource); err != nil {
		return NamedAPIResource{}, fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return resource, nil
}
