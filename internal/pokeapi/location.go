package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Location represents detailed information about a specific location area,
// including the Pokemon that can be encountered there.
type Location struct {
	// EncounterMethodRates describes the rates at which Pokemon can be encountered
	// using different methods (e.g., walking, surfing).
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	// GameIndex is the internal ID of this location area in the game.
	GameIndex int `json:"game_index"`
	// ID is the unique identifier for this location area resource.
	ID int `json:"id"`
	// Location contains the name and URL of the broader location this area belongs to.
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	// Name is the name of the location area.
	Name string `json:"name"`
	// Names lists the name of this location area in various languages.
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	// PokemonEncounters lists all Pokemon that can be found in this area,
	// along with details about how to encounter them.
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

// GetLocation returns the set of possible Pokémon encounters for an area.
func GetLocation(area string) (Location, error) {
	url := MapEndpoint + area

	var data []byte

	data, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return Location{}, fmt.Errorf("http.Get failed: %w", err)
		}

		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return Location{}, fmt.Errorf("ioutil.ReadAll failed: %w", err)
		}

		cache.Add(url, data)
	}

	var resource Location

	err := json.Unmarshal(data, &resource)
	if err != nil {
		return Location{}, fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return resource, nil
}
