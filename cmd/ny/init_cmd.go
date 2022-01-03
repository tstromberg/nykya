package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tstromberg/nykya/pkg/nykya"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type initCmd struct{}

func (a *initCmd) Run(globals *Globals) error {
	if err := os.MkdirAll(globals.Root, 0o700); err != nil {
		return err
	}

	dc, err := nykya.ConfigFromRoot(globals.Root)
	if err != nil {
		return fmt.Errorf("config from root: %w", err)
	}

	bs, err := yaml.Marshal(dc)
	if err != nil {
		return err
	}

	dest := filepath.Join(dc.Root, nykya.ConfigFileName)
	klog.Infof("Writing configuration to %s ...", dest)
	if err := os.WriteFile(dest, bs, 0o600); err != nil {
		return err
	}

	return nil
}
