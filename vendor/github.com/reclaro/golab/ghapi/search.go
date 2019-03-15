package ghapi

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
)

// Search executes a search inside a specified repositroy name
func Search(repoName string) (*github.CodeSearchResult, error) {
	searchString := fmt.Sprintf("TODO in:file repo:%s", repoName)
	client := github.NewClient(nil)
	ctx := context.Background()
	results, _, err := client.Search.Code(ctx, searchString, nil)
	if err != nil {
		return nil, err
	}
	return results, nil
}
