package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
)

func main() {
	// Replace these values with your CodeCommit repository details
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	repositoryURL := "https://git-codecommit.us-east-1.amazonaws.com/v1/repos/kellie1"
	userName := os.Getenv("GIT_USERNAME")
	accessToken := os.Getenv("GIT_ACCESS_TOKEN")

	GoGitGetCommitsByRepository(repositoryURL, userName, accessToken)
}

// Uses go-git to print a list of commits by repository
func GoGitGetCommitsByRepository(repositoryURL string, userName string, accessToken string) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  repositoryURL,
		Auth: &http.BasicAuth{Username: userName, Password: accessToken},
	})
	CheckIfError(err)

	// ... retrieves the branch pointed by HEAD
	ref, err := repo.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	CheckIfError(err)
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
