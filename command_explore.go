package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("you must provide a location area name")
	}

	locationResp, err := cfg.pokeapiClient.ListPokemon(args[0])

	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationResp.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}
