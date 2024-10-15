package logic

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"bitbucket-cli/internal/bitbucket"
	"bitbucket-cli/internal/storage"
)

const maxConcurrentClones = 10

func CloneRepositories(client *bitbucket.Client, projectKey, directory string, db *storage.SQLiteDB) {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	repos, err := client.FetchRepositories(projectKey)
	if err != nil {
		log.Fatalf("Failed to fetch repositories: %v", err)
	}

	oldRepos, err := db.GetOldRepositories()
	if err != nil {
		log.Fatalf("Failed to get old repositories: %v", err)
	}

	repoMap := make(map[string]bitbucket.Repository)
	for _, repo := range repos {
		repoMap[repo.Links.Clone[0].Href] = repo
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentClones)

	for _, oldRepo := range oldRepos {
		if repo, exists := repoMap[oldRepo.URL]; exists {
			wg.Add(1)
			sem <- struct{}{}

			go func(repo bitbucket.Repository) {
				defer wg.Done()
				defer func() { <-sem }()

				repoPath := filepath.Join(directory, repo.Name)
				cloneURL := repo.Links.Clone[0].Href

				if _, err := os.Stat(repoPath); os.IsNotExist(err) {
					fmt.Printf("Cloning repository: %s\n", repo.Name)
					if err := gitClone(cloneURL, directory); err != nil {
						log.Printf("Failed to clone repository %s: %v", repo.Name, err)
					}
				} else {
					fmt.Printf("Updating repository: %s\n", repo.Name)
					if err := gitFetchAndPull(repoPath); err != nil {
						log.Printf("Failed to update repository %s: %v", repo.Name, err)
					}
				}

				if err := db.UpdateRepository(cloneURL); err != nil {
					log.Printf("Failed to update database for repository %s: %v", repo.Name, err)
				}
			}(repo)
		}
	}

	wg.Wait()
}

func gitClone(url, directory string) error {
	cmd := exec.Command("git", "clone", url)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitFetchAndPull(repoPath string) error {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "pull")
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
