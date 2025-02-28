package main

import "math/rand"

func pokemanCaught(baseXP int) bool {
	divisor := float32(2)
	threshold := 100 - float32(baseXP)/divisor
	randomNumber := float32(rand.Intn(101))
	return randomNumber <= threshold
}
