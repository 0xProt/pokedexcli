package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PokeGet(pageUrl *string) ([]byte, error) {
	url := baseURL + "/location-area"
	if pageUrl != nil {
		url = *pageUrl
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and \nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

func PokeUnmarshal(body []byte) (PokeLocation, error) {
	location := PokeLocation{}
	err := json.Unmarshal(body, &location)
	if err != nil {
		return PokeLocation{}, err
	}
	return location, nil
}
