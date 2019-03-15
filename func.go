package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	fdk "github.com/fnproject/fdk-go"
	"github.com/reclaro/golab/ghapi"
	"io"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// you can invoke it passing the repository name via the Fn cli
	// echo -n '{"repo_name": "fnproject/fn", "query": "TODO"}'| fn invoke oracle-code fngh | jq .
	sq := &ghapi.SearchQuery{}
	err := json.NewDecoder(in).Decode(sq)
	if err != nil && err != io.EOF {
		_ = json.NewEncoder(out).Encode(&ghapi.SearchResults{Error: fmt.Sprintf("Error in decoding input %s", err.Error())})
		return
	}
	if sq.Repo == "" {
		sq.Repo = "fnproject/fn"
	}
	if sq.Query == "" {
		sq.Query = "TODO"
	}
	msg := ghapi.Search(sq)
	err = json.NewEncoder(out).Encode(&msg)
	if err != nil {
		_ = json.NewEncoder(out).Encode(&ghapi.SearchResults{Error: errors.New("Unable to encode results").Error()})
		return
	}
}
