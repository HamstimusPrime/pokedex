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
	config.pokedex = &map[string]pokemon{}
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
	pokedex    *map[string]pokemon
}

type pokeResult struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
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
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Throwing a ball to a provided pokemon that may or may not catch it",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "catch",
			description: "Verify if a certain pokemon has been caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "catch",
			description: "Prints a list of all capture",
			callback:    commandPokedex,
		},
	}
}

func commandExit(config *configuration, args ...string) error {
	if args != nil {
		fmt.Printf("'%v' command does not accept arguments\n", args[0])
		return nil
	}
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *configuration, args ...string) error {
	if args != nil {
		fmt.Printf("'%v' command does not accept arguments\n", args[0])
		return nil
	}
	allCommands := getCommands()
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range allCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(config *configuration, args ...string) error {
	if args != nil {
		fmt.Printf("'%v' command does not accept arguments\n", args[0])
		return nil
	}
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
	if args != nil {
		fmt.Printf("'%v' command does not accept arguments\n", args[0])
		return nil
	}
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

func commandExplore(config *configuration, location ...string) error {
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

func commandCatch(config *configuration, pkm ...string) error {

	if pkm == nil || len(pkm) > 1 {
		fmt.Println("no arguments provided\nThe catch command requires a valid pokemon name.")
		return nil
	}

	if _, exist := (*config.pokedex)[pkm[0]]; exist {
		fmt.Printf("%v already caught\n", pkm[0])
		return nil
	}

	apiURL := "https://pokeapi.co/api/v2/pokemon/"
	requestURL := apiURL + pkm[0]
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
	var newPokemon pokemon
	if err := decoder.Decode(&newPokemon); err != nil {
		fmt.Printf("error decoding request body, error: %v", err)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", newPokemon.Name)
	if pokemanCaught(newPokemon.BaseExperience) {
		fmt.Printf("%v was caught!\n", newPokemon.Name)
		(*config.pokedex)[newPokemon.Name] = newPokemon
		return nil
	} else {
		fmt.Printf("%v escaped!\n", newPokemon.Name)
	}
	return nil

}

func commandInspect(config *configuration, pkm ...string) error {
	if pkm == nil {
		fmt.Println("please provide a pokemon name")
		return nil
	}
	_, exists := (*config.pokedex)[pkm[0]]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	} else {
		fmt.Printf("Name: %v\n", (*config.pokedex)[pkm[0]].Name)
		fmt.Printf("Height: %v\n", (*config.pokedex)[pkm[0]].Height)
		fmt.Printf("Weight: %v\n", (*config.pokedex)[pkm[0]].Weight)
		fmt.Println("Stats:")
		for _, stat := range (*config.pokedex)[pkm[0]].Stats {
			fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Println("Types:")
		for _, stat := range (*config.pokedex)[pkm[0]].Types {
			fmt.Printf("  - %v\n", stat.Type.Name)
		}
	}
	return nil
}

func commandPokedex(config *configuration, args ...string) error {
	if args != nil {
		fmt.Printf("'%v' command does not accept arguments\n", args[0])
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range *config.pokedex {
		fmt.Printf(" - %v\n", pokemon.Name)
	}
	return nil
}
