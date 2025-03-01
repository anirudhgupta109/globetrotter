package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

// Destination represents the structure of each record in dataset.json.
type destination struct {
	City    string   `json:"city"`
	Country string   `json:"country"`
	Clues   []string `json:"clues"`
	FunFact []string `json:"fun_fact"`
	Trivia  []string `json:"trivia"`
}

func MigrateDestinations() {
	// Open the JSON dataset file.
	file, err := os.Open("schemas/dataset.json")
	if err != nil {
		log.Fatalf("Error opening dataset.json: %v", err)
	}
	defer file.Close()

	// Read the file content.
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading dataset.json: %v", err)
	}

	// Unmarshal the JSON into a slice of Destination.
	var destinations []destination
	if err := json.Unmarshal(data, &destinations); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Process each destination.
	for _, dest := range destinations {
		// Insert into destinations table and get the generated id.
		var destID uuid.UUID
		err = db.QueryRow(context.Background(),
			`INSERT INTO destinations (city, country) VALUES ($1, $2) RETURNING id`,
			&dest.City, &dest.Country,
		).Scan(&destID)
		if err != nil {
			log.Printf("Error inserting destination (%s, %s): %v", dest.City, dest.Country, err)
			continue
		}
		fmt.Printf("Inserted destination %s with id %d\n", dest.City, destID)

		// Insert clues.
		for _, clue := range dest.Clues {
			_, err = db.Exec(context.Background(),
				`INSERT INTO clues (destination_id, clue_text) VALUES ($1, $2)`,
				&destID, &clue,
			)
			if err != nil {
				log.Printf("Error inserting clue for destination id %d: %v", destID, err)
			}
		}

		// Insert fun facts.
		for _, fact := range dest.FunFact {
			_, err = db.Exec(context.Background(),
				`INSERT INTO fun_facts (destination_id, fact_text) VALUES ($1, $2)`,
				&destID, &fact,
			)
			if err != nil {
				log.Printf("Error inserting fun fact for destination id %d: %v", destID, err)
			}
		}

		// Insert trivia.
		for _, trivia := range dest.Trivia {
			_, err = db.Exec(context.Background(),
				`INSERT INTO trivia (destination_id, trivia_text) VALUES ($1, $2)`,
				&destID, &trivia,
			)
			if err != nil {
				log.Printf("Error inserting trivia for destination id %d: %v", destID, err)
			}
		}
	}

	fmt.Println("Data insertion complete!")
}
