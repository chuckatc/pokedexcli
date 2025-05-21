package main

import (
	"bufio"
	"fmt"
	"log"
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
	}

	config := cmdConfig{
		cache:       pokecache.NewCache(5 * time.Second),
		cmdRegistry: cmdRegistry,
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
		if ok {
			err := cliCmd.callback(&config, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
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
