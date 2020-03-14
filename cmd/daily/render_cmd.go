package main

import (
	"context"
	"fmt"

	"github.com/tstromberg/daily/pkg/action"
	"github.com/tstromberg/daily/pkg/daily"
	"k8s.io/klog"
)

type RenderCmd struct{}

func (c *RenderCmd) Run(globals *Globals) error {
	dc, err := daily.ConfigFromRoot(globals.Root)
	if err != nil {
		return fmt.Errorf("config from root: %w", err)
	}
	paths, err := action.Render(context.Background(), dc)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	klog.Infof("rendered paths: %v", paths)
	return nil
}
