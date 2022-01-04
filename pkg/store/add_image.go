package store

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tstromberg/nykya/pkg/nykya"
	"k8s.io/klog/v2"
)

// addImage adds an image
func addImage(ctx context.Context, dc nykya.Config, opts AddOptions) error {
	klog.Infof("addImage %+v", opts)

	format := opts.Format
	if format == "" {
		format = formatForPath(opts.Content)
	}

	i := nykya.RenderInput{
		FrontMatter: nykya.FrontMatter{
			Kind:   opts.Kind,
			Date:   nykya.NewYAMLTime(opts.Timestamp),
			Origin: opts.Content,
		},
		Format: formatForPath(opts.Content),
	}

	outPath, err := localPath(dc, i.FrontMatter)
	if err != nil {
		return fmt.Errorf("local path: %w", err)
	}

	return saveItem(ctx, dc, i, filepath.Join(dc.In, outPath))
}
