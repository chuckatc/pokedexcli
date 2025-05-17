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
	callback    func(*cmdConfig) error
}

var cmdRegistry map[string]cliCommand

type cmdConfig struct {
	cache    *pokecache.Cache
	Next     string
	Previous string
}

func main() {
	cmdRegistry = map[string]cliCommand{
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
	}

	repl()
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)

	config := cmdConfig{}
	config.cache = pokecache.NewCache(5 * time.Second)

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

		cliCmd, ok := cmdRegistry[command]
		if ok {
			err := cliCmd.callback(&config)
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

func commandHelp(config *cmdConfig) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range cmdRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit(config *cmdConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *cmdConfig) error {
	mapData := pokeapi.GetMap(config.Next, config.cache)
	config.Next = mapData.Next
	config.Previous = mapData.Previous
	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapB(config *cmdConfig) error {
	mapData := pokeapi.GetMap(config.Previous, config.cache)
	config.Next = mapData.Next
	config.Previous = mapData.Previous
	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}
