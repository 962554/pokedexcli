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

var (
	MapEndpoint string           = "https://pokeapi.co/api/v2/location-area/"
	cache       *pokecache.Cache = pokecache.NewCache(interval)
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

	var data []byte
	data, ok := cache.Get(url)
	if !ok {
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
