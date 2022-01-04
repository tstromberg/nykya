package store

import (
	"context"

	"github.com/tstromberg/nykya/pkg/nykya"
	"k8s.io/klog/v2"
)

// Scan scans a directory for posted content
func Scan(ctx context.Context, root string) ([]*nykya.RenderInput, error) {
	klog.Infof("Scanning root %s ...", root)
	return fromDirectory(root, root)
}
