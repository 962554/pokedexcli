package pokecache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/962554/pokedexcli/internal/pokecache"
)

func TestAddGet(t *testing.T) {
	t.Parallel()

	const interval = 5 * time.Second

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}
	for i, tcase := range cases {
		t.Run(fmt.Sprintf("Test tcase %v", i), func(t *testing.T) {
			t.Parallel()

			cache := pokecache.NewCache(interval)
			cache.Add(tcase.key, tcase.val)
			val, ok := cache.Get(tcase.key)

			if !ok {
				t.Errorf("expected to find key")

				return
			}

			if string(val) != string(tcase.val) {
				t.Errorf("expected to find value")

				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	t.Parallel()

	const (
		baseTime = 5 * time.Millisecond
		waitTime = baseTime + 5*time.Millisecond
	)

	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, hit := cache.Get("https://example.com")
	if !hit {
		t.Errorf("expected to find key")

		return
	}

	time.Sleep(waitTime)

	_, hit = cache.Get("https://example.com")
	if hit {
		t.Errorf("expected to not find key")

		return
	}
}
