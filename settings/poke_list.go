package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListPokemon(locationName string) (RespLocationArea, error) {
	url := baseURL + "/location-area/" + locationName

	if data, ok := c.cache.Get(url); ok {
		locationResp := RespLocationArea{}
		err := json.Unmarshal(data, &locationResp)
		if err != nil {
			return RespLocationArea{}, err

		}
		return locationResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespLocationArea{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespLocationArea{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespLocationArea{}, err
	}
	c.cache.Add(url, dat)
	locationResp := RespLocationArea{}
	err = json.Unmarshal(dat, &locationResp)
	if err != nil {
		return RespLocationArea{}, err
	}

	return locationResp, nil

}

func (c *Client) GetPokemon(pokemon string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemon

	if data, ok := c.cache.Get(url); ok {
		pokemonResp := Pokemon{}
		err := json.Unmarshal(data, &pokemonResp)
		if err != nil {
			return Pokemon{}, err

		}
		return pokemonResp, nil
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

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}
	c.cache.Add(url, dat)
	pokemonResp := Pokemon{}
	err = json.Unmarshal(dat, &pokemonResp)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemonResp, nil

}
