package main

import (
	"fmt"
	"os"

	"github.com/0xProt/pokedexcli/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func commandExit(cfg *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *Config) error {
	body, err := pokeapi.PokeGet(cfg.nextLocationsUrl)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}
	locations, err := pokeapi.PokeUnmarshal(body)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}
	cfg.nextLocationsUrl = locations.Next
	if locations.Previous == nil {
		cfg.prevLocationsUrl = nil
	} else {
		cfg.prevLocationsUrl = locations.Previous
	}
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return nil
}

func commandMapBack(cfg *Config) error {
	if *cfg.prevLocationsUrl == "" {
		fmt.Println("You're on the first page of locations")
		return nil
	}
	body, err := pokeapi.PokeGet(cfg.prevLocationsUrl)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}
	locations, err := pokeapi.PokeUnmarshal(body)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}
	cfg.nextLocationsUrl = locations.Next
	if locations.Previous == nil {
		cfg.prevLocationsUrl = nil
	} else {
		cfg.prevLocationsUrl = locations.Previous
	}
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
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
	}
}
