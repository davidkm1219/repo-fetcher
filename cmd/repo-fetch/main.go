package main

import (
	"log"
	"os"

	"bitbucket-cli/internal/bitbucket"
	"bitbucket-cli/internal/logic"
	"bitbucket-cli/internal/storage"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: <program> <project-key> <directory> <db-path>")
	}

	projectKey := os.Args[1]
	directory := os.Args[2]
	dbPath := os.Args[3]

	apiToken := os.Getenv("BITBUCKET_API_TOKEN")
	if apiToken == "" {
		log.Fatal("BITBUCKET_API_TOKEN environment variable is not set")
	}

	client := bitbucket.NewClient(apiToken)

	db, err := storage.NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	logic.CloneRepositories(client, projectKey, directory, db)
}
