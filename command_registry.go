package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
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
	fmt.Printf("Exploring %s...\n", args)
	pokemonList, err := cfg.pokeapiClient.PokeGetPokemon(args)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemonList.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *Config, args string) error {
	if args == "" {
		return errors.New("you must provide the name of a pokemon to catch")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args)
	pokeStats, err := cfg.pokeapiClient.PokeCatchPokemon(args)
	if err != nil {
		return err
	}
	catchChance := rand.IntN(500)
	if pokeStats.BaseExperience <= catchChance {
		fmt.Printf("%s was caught!\n", args)
		fmt.Print("You may now inspect it with the inspect command.\n")
		cfg.caughtPokemon[pokeStats.Name] = pokeStats
	} else {
		fmt.Printf("%s escaped!\n", args)
	}
	return nil
}

func commandInspect(cfg *Config, args string) error {
	if args == "" {
		return errors.New("you must provide the name of a pokemon to inspect")
	}

	tar, ok := cfg.caughtPokemon[args]
	if !ok {
		fmt.Print("you have not caught that pokemon\n")
	} else {

		fmt.Printf("Name: %s\n", tar.Name)
		fmt.Printf("Height: %d\n", tar.Height)
		fmt.Print("Stats:\n")
		for _, statEntry := range tar.Stats {
			fmt.Printf("  -%s: %d\n", statEntry.Stat.Name, statEntry.BaseStat)
		}
		fmt.Print("Types:\n")
		for _, typeEntry := range tar.Types {
			fmt.Printf("  - %s\n", typeEntry.Type.Name)
		}
	}
	return nil
}

func commandPokedex(cfg *Config, args string) error {
	if args != "" {
		return errors.New("pokedex doesn't need arguments, did you mean inspect?")
	}

	if len(cfg.caughtPokemon) == 0 {
		fmt.Println("your pokedex is currently empty, catch pokemon using the catch command")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokedexEntry := range cfg.caughtPokemon {
		fmt.Printf("  - %s\n", pokedexEntry.Name)
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch a specific Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "See stats for a specific caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Check your pokedex to see your caught Pokemon",
			callback:    commandPokedex,
		},
	}
}
