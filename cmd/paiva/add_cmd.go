package main

import (
	"context"
	"time"

	"github.com/tstromberg/paivalehti/pkg/paivalehti"
	"github.com/tstromberg/paivalehti/pkg/store"
)

type addCmd struct {
	Kind    string `arg required help:"Type of object to add"`
	Content string `arg required help:"Content to add"`

	Title string `help:"Set a title for the post"`
}

func (a *addCmd) Run(globals *Globals) error {
	dc, err := paivalehti.ConfigFromRoot(globals.Root)
	if err != nil {
		return err
	}

	return store.Add(context.Background(), dc, store.AddOptions{
		Content:   a.Content,
		Root:      globals.Root,
		Title:     a.Title,
		Kind:      a.Kind,
		Timestamp: time.Now(),
	})
}
