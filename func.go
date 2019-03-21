package main

import (
	"context"
	"encoding/json"
	"errors"

	fdk "github.com/fnproject/fdk-go"
	"github.com/reclaro/golab/ghapi"
	"io"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	sq := &ghapi.SearchQuery{Repo: "fnproject/fn", Query: "TODO"}
	msg := ghapi.Search(sq)
	err := json.NewEncoder(out).Encode(&msg)
	if err != nil {
		_ = json.NewEncoder(out).Encode(&ghapi.SearchResults{Error: errors.New("Unable to encode results").Error()})
		return
	}
}
