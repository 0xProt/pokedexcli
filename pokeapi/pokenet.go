package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) PokeGetLocation(pageUrl *string) (PokeLocation, error) {
	url := baseURL + "/location-area"
	if pageUrl != nil {
		url = *pageUrl
	}

	if val, ok := c.cache.Get(url); ok {
		location := PokeLocation{}
		err := json.Unmarshal(val, &location)
		if err != nil {
			return PokeLocation{}, err
		}
	}

	res, err := http.Get(url)
	if err != nil {
		return PokeLocation{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return PokeLocation{}, fmt.Errorf("response failed with status code: %d and \nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return PokeLocation{}, err
	}

	location := PokeLocation{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return PokeLocation{}, err
	}
	c.cache.Add(url, body)
	return location, nil
}

func (c *Client) PokeGetPokemon(locationName string) (WhichPokemonEncounters, error) {
	url := baseURL + "/location-area/" + locationName

	if val, ok := c.cache.Get(locationName); ok {
		whichPokemon := WhichPokemonEncounters{}
		err := json.Unmarshal(val, &whichPokemon)
		if err != nil {
			return WhichPokemonEncounters{}, err
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WhichPokemonEncounters{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return WhichPokemonEncounters{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return WhichPokemonEncounters{}, err
	}

	whichPokemon := WhichPokemonEncounters{}
	err = json.Unmarshal(data, &whichPokemon)
	if err != nil {
		return WhichPokemonEncounters{}, err
	}
	c.cache.Add(locationName, data)
	return whichPokemon, nil
}

func (c *Client) PokeCatchPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	if val, ok := c.cache.Get(pokemonName); ok {
		pokemonStats := Pokemon{}
		err := json.Unmarshal(val, &pokemonStats)
		if err != nil {
			return Pokemon{}, err
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemonStats := Pokemon{}
	err = json.Unmarshal(data, &pokemonStats)
	if err != nil {
		return Pokemon{}, err
	}
	c.cache.Add(pokemonName, data)
	return pokemonStats, nil
}
