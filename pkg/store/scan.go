package store

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/tstromberg/nykya/pkg/nykya"
	"k8s.io/klog"
)

// Scan scans a directory for posted content
func Scan(ctx context.Context, root string) ([]*nykya.RenderInput, error) {
	klog.Infof("Scanning root %s ...", root)

	fs, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*nykya.RenderInput

	for _, f := range fs {
		klog.V(1).Infof("Scanning subdir %s", f.Name())
		ds, err := fromDirectory(filepath.Join(root, f.Name()), root)
		if err != nil {
			return nil, fmt.Errorf("parse date: %w", err)
		}
		ps = append(ps, ds...)
	}
	return ps, nil
}
