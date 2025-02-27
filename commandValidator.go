package main

import "strings"

func IsValidCommand(input string) bool {
	inputList := strings.Split(input, " ")
	_, commandExists := getCommands()[inputList[0]]

	if commandExists && len(inputList) > 1 {
		return true
	}
	if commandExists && len(inputList) == 1 {
		return true
	}
	return false
}

func commandHasArguments(input string) (bool, string, string) {
	inputList := strings.Split(input, " ")

	if len(inputList) >= 2 {
		command := inputList[0]
		argument := inputList[1]
		return true, command, argument
	}

	return false, "", ""
}
