package ghapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

//SearchQuery is used to pass the repository name and the query we want to perform
type SearchQuery struct {
	Repo  string `json:"repo_name"`
	Query string `json:"query"`
}

type msgResult struct {
	Total    int    `json:"total_occurences,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
	HTMLURL  string `json:"html_url,omitempty"`
}

// SearchResults returns the results of the query
type SearchResults struct {
	Error   string       `json:"error"`
	Results []*msgResult `json:"results"`
}

// NewSearchQuery creates a new SearchQuery
func NewSearchQuery(reponame, searchTerm string) *SearchQuery {
	return &SearchQuery{Repo: reponame,
		Query: searchTerm}

}

func doSearch(searchTerm, repoName string) (*github.CodeSearchResult, error) {
	searchString := fmt.Sprintf("%s in:file repo:%s", searchTerm, repoName)
	client := github.NewClient(nil)
	ctx := context.Background()
	results, _, err := client.Search.Code(ctx, searchString, nil)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Search executes a search for a term inside a specified repositroy name
func Search(sq *SearchQuery) SearchResults {
	var msg SearchResults
	resu, err := doSearch(sq.Query, sq.Repo)
	if err != nil {
		msg = SearchResults{Error: fmt.Sprintf("No results found. %s ", err.Error())}
	} else {
		results := make([]*msgResult, 0)
		if len(resu.CodeResults) > 0 {
			recapResult := &msgResult{Total: *resu.Total}
			results = append(results, recapResult)
		}
		for _, v := range resu.CodeResults {
			r := &msgResult{
				FileName: *v.Name,
				FilePath: *v.Path,
				HTMLURL:  *v.HTMLURL,
			}
			results = append(results, r)
		}
		msg = SearchResults{Results: results}
	}
	return msg
}
