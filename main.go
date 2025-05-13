package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var cmdRegistry map[string]cliCommand

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
	}

	repl()
}

func repl() {
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

		cliCmd, ok := cmdRegistry[command]
		if ok {
			err := cliCmd.callback()
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

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range cmdRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
