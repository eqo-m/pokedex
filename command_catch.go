package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("please provide a pokemon name")
	}
	pokemonResp, err := cfg.pokeapiClient.GetPokemon(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonResp.Name)
	threshold := 50
	chance := rand.Intn(pokemonResp.BaseExperience)
	if chance > threshold {
		fmt.Printf("%s escaped!\n", pokemonResp.Name)
		return nil
	}

	fmt.Printf("%s was caught\n", pokemonResp.Name)
	cfg.caughtPokemon[args[0]] = pokemonResp
	return nil
}
