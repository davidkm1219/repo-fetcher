package main

import (
	"log"
	"os"

	"bitbucket-cli/internal/bitbucket"
	"bitbucket-cli/internal/logic"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: <program> <project-key> <directory> <api-token>")
	}

	projectKey := os.Args[1]
	directory := os.Args[2]
	apiToken := os.Args[3]

	client := bitbucket.NewClient(apiToken)
	logic.CloneRepositories(client, projectKey, directory)
}
