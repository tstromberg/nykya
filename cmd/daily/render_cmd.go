package main

import (
	"fmt"

	"github.com/tstromberg/daily/pkg/action"
	"github.com/tstromberg/daily/pkg/daily"
	"k8s.io/klog"
)

//
type RenderCmd struct{}

func renderCmd(root string) error {
	dc, err := daily.ConfigFromRoot(root)
	if err != nil {
		return fmt.Errorf("config from root: %w", err)
	}
	paths, err := action.Render(dc)
	if err != nil {
		return err
	}
	klog.Infof("rendered %d paths in %s", len(paths), dc.Out)
	return nil
}
