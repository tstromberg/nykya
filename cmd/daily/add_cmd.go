package main

import (
	"github.com/tstromberg/daily/pkg/action"
)

func addWithoutPath(root string, opts addOpts) error {
	return action.Add("", action.AddOptions{
		Root:        root,
		Description: opts.Description,
	})
}

func addPaths(root string, opts addOpts) error {
	for _, p := range opts.Paths {
		err := action.Add(p, action.AddOptions{
			Root:        root,
			Description: opts.Description,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
