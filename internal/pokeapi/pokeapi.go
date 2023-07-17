package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/neofight78/pokedexcli/internal/pokecache"
	"io"
	"net/http"
	"time"
)

type LocationAreasResult struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaResult struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type Client struct {
	cache pokecache.Cache
}

func NewClient() Client {
	return Client{
		cache: pokecache.NewCache(time.Minute * 5),
	}
}

func (client *Client) fetchUrl(url string) ([]byte, error) {
	cachedResult, ok := client.cache.Get(url)
	if ok {
		return cachedResult, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Response failed with status code: %d, and body:\n %s", res.StatusCode, body))
	}

	client.cache.Add(url, body)

	return body, nil
}

func (client *Client) FetchLocationAreas(url *string) (*LocationAreasResult, error) {
	var defaultUrl = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	if url == nil {
		url = &defaultUrl
	}

	body, err := client.fetchUrl(*url)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch url %s: %w", *url, err)
	}

	var result LocationAreasResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) FetchLocationArea(name string) (*LocationAreaResult, error) {
	body, err := client.fetchUrl(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name))
	if err != nil {
		return nil, fmt.Errorf("unable to fetch location %s: %w", name, err)
	}

	var result LocationAreaResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) FetchPokemon(name string) (*Pokemon, error) {
	body, err := client.fetchUrl(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name))
	if err != nil {
		return nil, fmt.Errorf("unable to fetch pokemon %s: %w", name, err)
	}

	var result Pokemon
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
