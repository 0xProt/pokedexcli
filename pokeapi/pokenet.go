package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) PokeGet(pageUrl *string) (PokeLocation, error) {
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
	return location, nil
}
