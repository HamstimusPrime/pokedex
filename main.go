package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	config := configuration{}
	for {
		fmt.Print("Pokedex> ")
		scanner.Scan()
		input := scanner.Text()

		if _, exists := getCommands()[input]; exists {
			inputCallBack := getCommands()[input].callback
			inputCallBack(&config)
		} else {
			fmt.Println("Unknown command")
		}

	}

}

func cleanInput(text string) []string {
	loweString := strings.Trim(strings.ToLower(text), " ")
	words := strings.Split(loweString, " ")
	return words
}

func commandExit(config *configuration) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *configuration) error {
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
	callback    func(*configuration) error
}

type configuration struct {
	Next     string
	Previous string
}

type pokeResult struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getCommands() map[string]cliCommand {

	var allCommands = map[string]cliCommand{
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
	}
	return allCommands
}

func commandMap(config *configuration) error {
	url := config.Next
	//====== Check if call to API has been initiated earlier ========
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	} else {
		url = config.Next
	}

	//================= Request logic===============//
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to access url %v", err)
	}
	defer res.Body.Close()
	//===== Parse response Body into Json ======//
	decoder := json.NewDecoder(res.Body)
	pokeMaps := pokeResult{}
	if err := decoder.Decode(&pokeMaps); err != nil {
		return fmt.Errorf("unable to decode response body to json: %v", err)
	}
	for _, v := range pokeMaps.Results {
		fmt.Println(v.Name)
	}
	config.Next = pokeMaps.Next
	config.Previous = pokeMaps.Previous
	return nil
}

func commandMapB(config *configuration) error {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	url := config.Previous
	//==============================================//
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to access url %v", err)
	}
	defer res.Body.Close()
	//===== Parse response Body into Json ======//
	decoder := json.NewDecoder(res.Body)
	pokeMaps := pokeResult{}
	if err := decoder.Decode(&pokeMaps); err != nil {
		return fmt.Errorf("unable to decode response body to json: %v", err)
	}
	for _, v := range pokeMaps.Results {
		fmt.Println(v.Name)
	}
	config.Previous = pokeMaps.Previous
	config.Next = pokeMaps.Next
	return nil
}
