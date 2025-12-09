package pokeapi

import (
	"net/http"
	"time"

	pokecache "github.com/eqo-m/pokedex/internal"
)

// Client -
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

// NewClient -
func NewClient(timeout time.Duration, cacheDuration time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheDuration),
	}
}
