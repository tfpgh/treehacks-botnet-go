package main

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"os"
)

func main() {
	botName, err := generateBotName()
	if err != nil {
		fmt.Printf("Failed to generate unique bot name, %v. Using \"%v\"\n", err, botName)
	} else {
		fmt.Printf("Generated unique bot name: \"%v\"\n", botName)
	}
}

func generateBotName() (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		// Failed to generate unique ID
		return "default", err
	}

	username := os.Getenv("USER")
	if len(username) > 0 {
		// Username was found. Returning username + id
		return fmt.Sprintf("%v-%v", username, id), nil
	} else {
		// No username. Just returning id
		return id, nil
	}

}
