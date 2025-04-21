package main

import (
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

func commandExit(cfg *Config, args string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *Config, args string) error {
	locations, err := cfg.pokeapiClient.PokeGetLocation(cfg.nextLocationsUrl)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}
	cfg.nextLocationsUrl = locations.Next
	cfg.prevLocationsUrl = locations.Previous

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return nil
}

func commandMapBack(cfg *Config, args string) error {
	if cfg.prevLocationsUrl == nil {
		return errors.New("you're on the first page of locations")
	}
	locations, err := cfg.pokeapiClient.PokeGetLocation(cfg.prevLocationsUrl)
	if err != nil {
		return err
	}
	cfg.nextLocationsUrl = locations.Next
	cfg.prevLocationsUrl = locations.Previous
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(cfg *Config, args string) error {
	if args == "" {
		return errors.New("explore requires a location area name as input")
	}
	pokemonList, err := cfg.pokeapiClient.PokeGetPokemon(args)
	if err != nil {
		return err
	}
	for _, pokemon := range pokemonList.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Shows the next 20 map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows the previous 20 map locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific area to see its Pokemon",
			callback:    commandExplore,
		},
	}
}
