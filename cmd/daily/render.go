package main

import (
	"github.com/tstromberg/daily/pkg/action"
	
	"github.com/urfave/cli/v2"
	"k8s.io/klog"
)

func renderCmd(c *cli.Context) error {
	src := c.Args().Get(0)
	// TODO: get from YAML
	dst := c.Args().Get(1)
	paths, err := action.Render(src, dst)
	if err != nil {
		return err
	}
	klog.Infof("rendered %d paths in %s", len(paths), dst)
	return nil
}
