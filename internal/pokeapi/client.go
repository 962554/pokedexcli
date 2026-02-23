package pokeapi

import (
	"net/http"
	"time"

	"github.com/962554/pokedexcli/internal/pokecache"
)

// Client is an HTTP client that communicates with the PokeAPI and caches results.
type Client struct {
	cache      *pokecache.Cache
	httpClient http.Client
}

// NewClient creates a new PokeAPI client with a cache that reaps entries after the given interval.
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
