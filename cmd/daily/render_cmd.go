package main

import (
	"github.com/tstromberg/daily/pkg/action"
	"github.com/tstromberg/daily/pkg/daily"
	"k8s.io/klog"
)

func renderCmd(root string) error {
	dc := daily.ConfigFromRoot(root)
	paths, err := action.Render(dc)
	if err != nil {
		return err
	}
	klog.Infof("rendered %d paths in %s", len(paths), dc.Out)
	return nil
}
