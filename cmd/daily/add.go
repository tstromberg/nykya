package main

import (
	"github.com/tstromberg/daily/pkg/action"
	
	"github.com/urfave/cli/v2"
)

func addCmd(c *cli.Context) error {
	path := c.Args().Get(0)
	return action.Add(path, rootFlag, action.AddOptions{})
}
