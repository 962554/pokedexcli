// Package pokeapi provides a client for interacting with the PokeAPI.
// It handles HTTP requests to various endpoints and manages data
// retrieval for location areas and other Pokemon-related resources.
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/962554/pokedexcli/internal/pokecache"
)

const interval = 60 * time.Second

// MapEndpoint is the default URL for the location-area resource.
var (
	MapEndpoint = "https://pokeapi.co/api/v2/location-area/"
	cache       = pokecache.NewCache(interval)
)

// NamedAPIResource represents a paginated response from the PokeAPI,
// containing metadata about the total count and links to adjacent pages.
type NamedAPIResource struct {
	// Count is the total number of resources available for this request.
	Count int `json:"count"`
	// Next is the URL for the next page of results, if any.
	Next *string `json:"next"`
	// Previous is the URL for the previous page of results, if any.
	Previous *string `json:"previous"`
	// Results contains the specific resource data for the current page.
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// GetLocationAreas returns all location-areas from PokeAPI
func GetLocationAreas(url string) (NamedAPIResource, error) {
	var data []byte
	data, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return NamedAPIResource{}, fmt.Errorf("http.Get failed: %w", err)
		}

		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return NamedAPIResource{}, fmt.Errorf("ioutil.ReadAll failed: %w", err)
		}
		cache.Add(url, data)
	}

	var resource NamedAPIResource

	if err := json.Unmarshal(data, &resource); err != nil {
		return NamedAPIResource{}, fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return resource, nil
}
