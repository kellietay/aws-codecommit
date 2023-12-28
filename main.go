package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
)

func main() {

	UseGoGit := false

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if UseGoGit {
		GoGitGetCommitsByRepository()
	}

	ctx := context.TODO()
	GetListRepos(ctx)

}

// Uses go-git to print a list of commits by repository
func GoGitGetCommitsByRepository() {
	repositoryURL := "https://git-codecommit.us-east-1.amazonaws.com/v1/repos/kellie1"
	userName := os.Getenv("GIT_USERNAME")
	accessToken := os.Getenv("GIT_ACCESS_TOKEN")

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

func GetAWSCodeCommitClient(ctx context.Context) (*codecommit.Client, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Create an Amazon S3 service client
	return codecommit.NewFromConfig(cfg), nil
}

func GetListRepos(ctx context.Context) {
	listRepositoriesInput := codecommit.ListRepositoriesInput{
		Order:  types.OrderEnumAscending,
		SortBy: types.SortByEnumModifiedDate,
	}
	client, err := GetAWSCodeCommitClient(ctx)
	CheckIfError(err)

	repositoryList, err := client.ListRepositories(ctx, &listRepositoriesInput)
	CheckIfError(err)

	for _, v := range repositoryList.Repositories {
		fmt.Printf("repository ID: %+v, repository Name: %+v \n", *v.RepositoryId, *v.RepositoryName)
		GetListBranches(ctx, v.RepositoryName)
	}
}

func GetListBranches(ctx context.Context, repositoryName *string) {
	listBranchesInput := codecommit.ListBranchesInput{
		RepositoryName: repositoryName,
	}

	client, err := GetAWSCodeCommitClient(ctx)
	CheckIfError(err)

	branchList, err := client.ListBranches(ctx, &listBranchesInput)
	CheckIfError(err)

	for _, b := range branchList.Branches {
		fmt.Printf("   --branch: %+v\n", b)
	}

}
