package main

import (
	"path/filepath"

	"github.com/tstromberg/daily/pkg/action"
	"k8s.io/klog"
)

func renderCmd(root string) error {
	src := filepath.Join(root, "in")
	dst := filepath.Join(root, "out")
	paths, err := action.Render(src, dst)
	if err != nil {
		return err
	}
	klog.Infof("rendered %d paths in %s", len(paths), dst)
	return nil
}
