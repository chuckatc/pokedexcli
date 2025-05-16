package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chuckatc/pokedexcli/internal/pokecache"
)

const LocationAreaUrl = "https://pokeapi.co/api/v2/location-area/"

type LocationAreaData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetMap(nextUrl string, cache *pokecache.Cache) LocationAreaData {
	var data LocationAreaData

	url := LocationAreaUrl
	if nextUrl != "" {
		url = nextUrl
	}

	mapData, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(mapData, &data)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("\t\tCACHE HIT")
		return data
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	cache.Add(url, body)

	// fmt.Println("\t\tCACHE MISS")
	return data
}
