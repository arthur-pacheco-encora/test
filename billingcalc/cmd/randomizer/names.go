package main

import (
	"errors"
	"fmt"
	"log"
)

var adjs = []string{"happy", "funny", "colorful", "cheerful", "mystic", "brave", "friendly", "enchanted", "epic", "resilient", "unique", "clever", "strong", "noble", "faithful", "serious", "tall", "short", "fast", "slow", "kind", "polite", "generous", "smooth", "shiny", "neat", "sweet", "healthy"}

var nouns = []string{"unicorn", "panda", "octopus", "dragon", "wizard", "phoenix", "robot", "mermaid", "centaur", "knight", "book", "bird", "worm", "whale", "tank", "fish", "coin", "apple", "door", "bag", "window", "television", "train", "plane", "phone", "tree", "computer"}

var numbers = []int{1, 2, 3, 4, 5}

var names = make(map[string]bool, len(adjs)*len(nouns))

func buildNames() {
	for _, adj := range adjs {
		for _, noun := range nouns {
			for _, num := range numbers {
				name := fmt.Sprintf("%s %s %d", adj, noun, num)
				names[name] = true
			}
		}
	}
	log.Printf("Total of available names: %d", len(names))
}

func getRandomName() (string, error) {
	var name string
	for k := range names {
		name = k
		break
	}

	if name == "" {
		return "", errors.New("Out of available names!")
	}

	delete(names, name)
	return name, nil
}
