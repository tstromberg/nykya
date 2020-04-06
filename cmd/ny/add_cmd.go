package main

import (
	"context"
	"time"

	"github.com/tstromberg/nykya/pkg/nykya"
	"github.com/tstromberg/nykya/pkg/store"
)

type addCmd struct {
	Kind    string `arg required help:"Type of object to add"`
	Content string `arg optional help:"Content to add"`
	Format  string `arg optional help:"Format of content (default is auto-detect)"`

	Title string `help:"Set a title for the post"`
}

func (a *addCmd) Run(globals *Globals) error {
	dc, err := nykya.ConfigFromRoot(globals.Root)
	if err != nil {
		return err
	}

	return store.Add(context.Background(), dc, store.AddOptions{
		Content:   a.Content,
		Root:      globals.Root,
		Title:     a.Title,
		Kind:      a.Kind,
		Format:    a.Format,
		Timestamp: time.Now(),
	})
}
