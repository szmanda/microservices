package main

import (
	"log"
)

func main() {
	err := LoadConfig("parameters.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v\n", appConfig)

	StartServer()
}
