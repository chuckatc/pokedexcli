package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chuckatc/pokedexcli/internal/pokeapi"
	"github.com/chuckatc/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*cmdConfig, []string) error
}

type cmdConfig struct {
	cache       *pokecache.Cache
	cmdRegistry map[string]cliCommand
	pokedex     map[string]pokeapi.PokemonData
	Next        string
	Previous    string
}

func main() {
	cmdRegistry := map[string]cliCommand{
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
			description: "Show next map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show previous map locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show Pokemon in your Pokedex",
			callback:    commandPokedex,
		},
	}

	config := cmdConfig{
		cache:       pokecache.NewCache(5 * time.Second),
		cmdRegistry: cmdRegistry,
		pokedex:     make(map[string]pokeapi.PokemonData),
	}

	repl(config)
}

func repl(config cmdConfig) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		input := scanner.Text()

		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}
		command := words[0]

		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		cliCmd, ok := config.cmdRegistry[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := cliCmd.callback(&config, args)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	textLower := strings.ToLower(text)
	words := strings.Fields(textLower)
	return words
}

func commandHelp(config *cmdConfig, args []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range config.cmdRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit(config *cmdConfig, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *cmdConfig, args []string) error {
	mapData := pokeapi.GetMap(config.Next, config.cache)
	config.Next = mapData.Next
	config.Previous = mapData.Previous
	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapB(config *cmdConfig, args []string) error {
	mapData := pokeapi.GetMap(config.Previous, config.cache)
	config.Next = mapData.Next
	config.Previous = mapData.Previous
	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *cmdConfig, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: explore <location_area>")
	}

	exploreData := pokeapi.GetExploreData(args[0], config.cache)

	fmt.Println("Found Pokemon:")
	for _, pokeEncounter := range exploreData.PokemonEncounters {
		fmt.Println("-", pokeEncounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *cmdConfig, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: catch <pokemon_name>")
	}
	name := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemonData, err := pokeapi.GetPokemonData(name, config.cache)
	if err != nil {
		return fmt.Errorf("you can't get ye %s", name)
	}

	if !attemptToCatch(pokemonData) {
		fmt.Println(name, "escaped!")
		return nil
	}

	fmt.Println(name, "was caught!")
	config.pokedex[name] = pokemonData

	return nil
}

func attemptToCatch(pokemonData pokeapi.PokemonData) bool {
	prob := probToCatch(pokemonData.BaseExperience)
	randFloat := rand.Float64()

	return randFloat > prob
}

func probToCatch(baseExp int) float64 {
	const maxBaseExp = 1000  // supposedly above maximum base experience of any pokemon
	const maxExpBuffer = 100 // in case there are some over maxBaseExp
	const probExponent = 1.2 // more base experience means exponentially harder to catch
	const probDivisor = 2    // divide 0-1 range to scale probability

	if baseExp > maxBaseExp-maxExpBuffer {
		baseExp = maxBaseExp - maxExpBuffer
	}

	prob := math.Pow(
		float64(maxBaseExp-baseExp)/maxBaseExp, probExponent) / probDivisor

	return prob
}

func commandInspect(config *cmdConfig, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: inspect <pokemon_name>")
	}
	name := args[0]

	pokemon, ok := config.pokedex[name]
	if !ok {
		return fmt.Errorf("you haven't caught %s yet", name)
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)

	}

	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokeType.Type.Name)
	}

	return nil
}

func commandPokedex(config *cmdConfig, args []string) error {
	fmt.Println("Your Pokedex:")

	for _, pokemon := range config.pokedex {
		fmt.Println("  -", pokemon.Name)
	}

	return nil
}
