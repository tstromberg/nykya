package main

import (
	"context"
	"time"

	"github.com/tstromberg/daily/pkg/action"
	"github.com/tstromberg/daily/pkg/daily"
)

type AddCmd struct {
	Kind    string `arg required help:"Type of object to add"`
	Content string `arg required help:"Content to add"`

	Title string `help:"Set a title for the post"`
}

func (a *AddCmd) Run(globals *Globals) error {
	dc, err := daily.ConfigFromRoot(globals.Root)
	if err != nil {
		return err
	}

	return action.Add(context.Background(), dc, action.AddOptions{
		Content:   a.Content,
		Root:      globals.Root,
		Title:     a.Title,
		Kind:      a.Kind,
		Timestamp: time.Now(),
	})
}
