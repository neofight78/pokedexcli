package main

import (
	"bufio"
	"fmt"
	"github.com/neofight78/pokedexcli/internal/pokeapi"
	"log"
	"math/rand"
	"os"
	"strings"
)

type config struct {
	previous *string
	next     *string
	client   *pokeapi.Client
	pokedex  map[string]*pokeapi.Pokemon
}

func newConfig() config {
	client := pokeapi.NewClient()

	return config{
		client:  &client,
		pokedex: make(map[string]*pokeapi.Pokemon),
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, params []string) error
}

func main() {
	config := newConfig()
	commands := commands()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		text := scanner.Text()
		parts := strings.Fields(text)

		if command, ok := commands[parts[0]]; ok {
			err := command.callback(&config, parts[1:])
			if err != nil {
				log.Fatalf("unhandled error: %s", err)
			}
		}
	}
}

func commandHelp(_ *config, _ []string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	commands := commands()

	for name := range commands {
		fmt.Printf("%s: %s\n", name, commands[name].description)
	}

	fmt.Println("")

	return nil
}

func commandExit(_ *config, _ []string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config, _ []string) error {
	if config.previous != nil && config.next == nil {
		fmt.Println("Cannot go forward any further")
		return nil
	}

	areas, err := config.client.FetchLocationAreas(config.next)
	if err != nil {
		return err
	}

	config.previous = areas.Previous
	config.next = areas.Next

	for _, location := range areas.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(config *config, _ []string) error {
	if config.previous == nil {
		fmt.Println("Cannot go back any further")
		return nil
	}

	areas, err := config.client.FetchLocationAreas(config.previous)
	if err != nil {
		return err
	}

	config.previous = areas.Previous
	config.next = areas.Next

	for _, location := range areas.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(config *config, parameters []string) error {
	area, err := config.client.FetchLocationArea(parameters[0])
	if err != nil {
		return err
	}

	for _, encounter := range area.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *config, parameters []string) error {
	pokemon, err := config.client.FetchPokemon(parameters[0])
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	if pokemon.BaseExperience < rand.Intn(200) {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		config.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(config *config, parameters []string) error {
	pokemon, ok := config.pokedex[parameters[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Printf("Types:\n")
	for _, ptype := range pokemon.Types {
		fmt.Printf("  - %s\n", ptype.Type.Name)
	}

	return nil
}

func commandPokedex(config *config, parameters []string) error {
	fmt.Printf("Your Pokedex:\n")
	for _, pokemon := range config.pokedex {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}

func commands() map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Explores locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Returns to previous locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explores the given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all captured pokemon",
			callback:    commandPokedex,
		},
	}
}
