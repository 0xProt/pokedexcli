package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/0xProt/pokedexcli/pokeapi"
)

type Config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsUrl *string
	prevLocationsUrl *string
	caughtPokemon    map[string]pokeapi.Pokemon
}

func startRepl(cfg *Config) {
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		argsValue := ""
		if len(words) > 1 {
			argsValue = words[1]
		}
		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, argsValue)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
