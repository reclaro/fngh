package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/go-github/github"
	//"github.com/google/go-github/v24/github"
)

func main() {
	//fdk.Handle(fdk.HandlerFunc(myHandler))
	GitHubCalls("")
}

//Item is a result item from code search
/* type Item struct {
	Name string `json:"name"`
	path string `json:"path"`
}

// CodeSearchResult contains results from a search in the code
type CodeSearchResult struct {
	Repository string
	Tot        int    `json: "total_count"`
	Items      []Item `json:"items"`
} */

func GitHubCalls(repo string) (*github.CodeSearchResult, error) {
	if repo == "" {
		repo = "fnproject/fn"
	}
	searchString := fmt.Sprintf("TODO in:file repo:%s", repo)
	client := github.NewClient(nil)
	ctx := context.Background()
	results, _, err := client.Search.Code(ctx, searchString, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Total occurences: %d\n\n", *results.Total)
	/* 	for _, cr := range results.CodeResults {
		fmt.Printf("File name %s\n", cr.GetName())
		fmt.Printf("File path %s\n", cr.GetPath())
		fmt.Printf("Html URL %s\n\n", cr.GetHTMLURL())

	} */
	return results, nil

}

type SearchQuery struct {
	Repo string `json:"repo_name"`
}

type SearchResults struct {
	Error string                   `json:"error"`
	Msg   *github.CodeSearchResult `json:"results"`
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	p := &SearchQuery{Repo: "World"}
	json.NewDecoder(in).Decode(p)
	resu, err := GitHubCalls(p.Repo)
	var msg SearchResults
	if err != nil {
		msg = SearchResults{Error: err.Error()}
	} else {
		msg = SearchResults{Msg: resu}
	}
	/* msg := struct {
		Msg string `json:"message"`
	}{
		Msg: fmt.Sprintf("Hello %s", p.Repo),
	}
	json.NewEncoder(out).Encode(&msg) */
	json.NewEncoder(out).Encode(&msg)
}
