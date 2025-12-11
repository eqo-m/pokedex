package main

import (
	"fmt"
)

func commandPokedex(cfg *config, args []string) error {

	fmt.Println("Your Pokemon:")

	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}

	return nil
}
