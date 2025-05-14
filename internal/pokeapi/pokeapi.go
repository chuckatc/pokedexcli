package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func GetMap(nextUrl string) LocationAreaData {
	url := LocationAreaUrl
	if nextUrl != "" {
		url = nextUrl
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

	var data LocationAreaData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
