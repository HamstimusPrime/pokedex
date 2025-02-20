package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex> ")
		scanner.Scan()
		input := scanner.Text()

		if _, exists := getCommands()[input]; exists {
			inputCallBack := getCommands()[input].callback
			inputCallBack()
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

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
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
	callback    func() error
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
	}
	return allCommands
}
