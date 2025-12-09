package main

import (
	"time"

	pokeapi "github.com/eqo-m/pokedex/settings"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
	}

	startRepl(cfg)
}
