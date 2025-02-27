package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedex/internal/pokecache"
	"strings"
	"time"
)

func main() {
	apiURL := "https://pokeapi.co/api/v2/location-area?limit=20"
	scanner := bufio.NewScanner(os.Stdin)
	config := configuration{}
	config.pokeApiURL = apiURL
	config.cache = pokecache.NewCache(5 * time.Second)
	config.apiResults = &pokeResult{}
	config.location = &pokeLocation{}
	for {
		fmt.Print("Pokedex> ")
		scanner.Scan()
		input := scanner.Text()
		if IsValidCommand(input) {
			hasArgument, command, argument := commandHasArguments(input)
			if hasArgument {
				getCommands()[command].callback(&config, argument)
				continue
			} else {
				getCommands()[input].callback(&config)
				continue
			}
		}
		fmt.Println("Unknown command")

	}

}

func cleanInput(text string) []string {
	loweString := strings.Trim(strings.ToLower(text), " ")
	words := strings.Split(loweString, " ")
	return words
}

func commandExit(config *configuration, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *configuration, args ...string) error {
	allCommands := getCommands()
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range allCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*configuration, ...string) error
}

type configuration struct {
	pokeApiURL string
	Next       string
	Previous   string
	cache      *pokecache.Cache
	apiResults *pokeResult
	location   *pokeLocation
}

type pokeResult struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokeLocation struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func getCommands() map[string]cliCommand {

	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Return pokemon Maps",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Returns the previous pokemon Maps",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Find various pokemons that can be found in a given location",
			callback:    explore,
		},
	}
}

func commandMap(config *configuration, args ...string) error {
	var data []byte

	if config.pokeApiURL == "" {
		return fmt.Errorf("you are on the last page")
	}

	byteData, exist := config.cache.Get(config.pokeApiURL)

	if !exist {
		res, err := http.Get(config.pokeApiURL)
		if err != nil {
			return fmt.Errorf("could not make API request, error: %v", err)
		}
		defer res.Body.Close()

		d, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		config.cache.Add(config.pokeApiURL, d)

		if err := json.Unmarshal(d, &config.apiResults); err != nil {
			return fmt.Errorf("unable to parse Json, error: %v", err)
		}
		for _, result := range config.apiResults.Results {
			fmt.Println(result.Name)
		}
		config.pokeApiURL = config.apiResults.Next
		config.Next = config.apiResults.Next
		config.Previous = config.apiResults.Previous
		return nil
	} else {
		data = byteData
		if err := json.Unmarshal(data, config.apiResults); err != nil {
			return fmt.Errorf("unable to parse Json, error: %v", err)
		}
		for _, results := range config.apiResults.Results {
			fmt.Println(results.Name)
		}
		config.pokeApiURL = config.Next
		config.Previous = config.apiResults.Previous
		return nil
	}

}

func commandMapB(config *configuration, args ...string) error {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	byteData, exists := config.cache.Get(config.Previous)
	if exists {
		if err := json.Unmarshal(byteData, config.apiResults); err != nil {
			return fmt.Errorf("unable to parse Json, error: %v", err)
		}
		for _, results := range config.apiResults.Results {
			fmt.Println(results.Name)
		}
	} else {

	}
	config.Previous = config.apiResults.Previous

	return nil
}

func explore(config *configuration, location ...string) error {
	if location == nil || len(location) > 1 {
		fmt.Println("no arguments provided\nThe explore command requires a valid location.")
		return nil
	}
	apiURL := "https://pokeapi.co/api/v2/location-area/"
	requestURL := apiURL + location[0]
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("unable to make request err: %v", err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode > 200 {
		fmt.Printf("request error! status code: %v\n", res.StatusCode)
		return nil
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&config.location); err != nil {
		fmt.Printf("error decoding request body, error: %v", err)
		return nil

	}
	fmt.Println("Exploring pastoria-city-area...\nFound Pokemon:")
	for _, pokeEnc := range config.location.PokemonEncounters {
		fmt.Println("- " + pokeEnc.Pokemon.Name)
	}
	return nil
}
