package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	fdk "github.com/fnproject/fdk-go"
	"github.com/reclaro/golab/ghapi"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type SearchQuery struct {
	Repo string `json:"repo_name"`
}

type SearchResults struct {
	Error   string       `json:"error"`
	Results []*msgResult `json:"results"`
}

type msgResult struct {
	Total    int    `json:"total_occurences,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
	HTMLURL  string `json:"html_url,omitempty"`
}

func searchInRepo(repo string) SearchResults {
	var msg SearchResults
	resu, err := ghapi.Search(repo)
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

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// you can invoke it passing the repository name via the Fn cli
	// echo -n '{"repo_name": "fnproject/fn"}'| fn invoke oracle-code fngh | jq .
	p := &SearchQuery{}
	err := json.NewDecoder(in).Decode(p)
	if err != nil && err != io.EOF {
		_ = json.NewEncoder(out).Encode(&SearchResults{Error: fmt.Sprintf("Error in decoding input %s", err.Error())})
		return
	}
	if p.Repo == "" {
		p.Repo = "fnproject/fn"
	}
	msg := searchInRepo(p.Repo)
	err = json.NewEncoder(out).Encode(&msg)
	if err != nil {
		_ = json.NewEncoder(out).Encode(&SearchResults{Error: errors.New("Unable to encode results").Error()})
		return
	}
}
