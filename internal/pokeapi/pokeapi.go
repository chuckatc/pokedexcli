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

type LocationAreaDetailData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetExploreData(locationArea string, cache *pokecache.Cache) LocationAreaDetailData {
	var data LocationAreaDetailData
	url := LocationAreaUrl + locationArea

	exploreData, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(exploreData, &data)
		if err != nil {
			log.Fatal(err)
		}
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

	return data
}
