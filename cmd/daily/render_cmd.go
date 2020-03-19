package main

import (
	"context"
	"fmt"

	"github.com/tstromberg/daily/pkg/daily"
	"github.com/tstromberg/daily/pkg/render"
	"github.com/tstromberg/daily/pkg/store"
	"k8s.io/klog"
)

type renderCmd struct{}

func (c *renderCmd) Run(globals *Globals) error {
	dc, err := daily.ConfigFromRoot(globals.Root)
	if err != nil {
		return fmt.Errorf("config from root: %w", err)
	}
	ctx := context.Background()
	items, err := store.Scan(ctx, dc.Root)
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
