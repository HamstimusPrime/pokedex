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
		if input == "exit" {
			return
		}

	}

}

func cleanInput(text string) []string {
	loweString := strings.Trim(strings.ToLower(text), " ")
	stringSlice := strings.Split(loweString, " ")
	return stringSlice
}
