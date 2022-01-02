package main

import (
	"context"
	"fmt"

	"github.com/tstromberg/nykya/pkg/nykya"
	"github.com/tstromberg/nykya/pkg/render"
	"github.com/tstromberg/nykya/pkg/store"
	"k8s.io/klog/v2"
)

type renderCmd struct{}

func (c *renderCmd) Run(globals *Globals) error {
	dc, err := nykya.ConfigFromRoot(globals.Root)
	if err != nil {
		return fmt.Errorf("config from root: %w", err)
	}
	ctx := context.Background()
	items, err := store.Scan(ctx, dc.In)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	paths, err := render.Site(ctx, dc, items)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	klog.Infof("rendered paths: %v", paths)
	return nil
}
