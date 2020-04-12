package main

import (
	"context"
	"time"

	"github.com/tstromberg/nykya/pkg/nykya"
	"github.com/tstromberg/nykya/pkg/store"
)

type addCmd struct {
	Kind    string   `arg required help:"Type of object to add"`
	Title   string   `optional help:"Title of object"`
	Format  string   `optional help:"Format of content (default is auto-detect)"`
	Content []string `arg name:"content" help:"Content to add"`
}

func (a *addCmd) Run(globals *Globals) error {
	dc, err := nykya.ConfigFromRoot(globals.Root)
	if err != nil {
		return err
	}

	for _, c := range a.Content {
		err := store.Add(context.Background(), dc, store.AddOptions{
			Content:   c,
			Title:     a.Title,
			Kind:      a.Kind,
			Format:    a.Format,
			Timestamp: time.Now(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
