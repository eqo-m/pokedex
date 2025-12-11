package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, args []string) error {

	if len(args) == 0 {
		return errors.New("please provide a pokemon name")
	}

	name := args[0]

	pokemon, ok := cfg.caughtPokemon[name]
	if !ok {
		return errors.New("haven't caught that one yet")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s : %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}
