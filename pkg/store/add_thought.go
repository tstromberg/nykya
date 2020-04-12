package store

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tstromberg/nykya/pkg/nykya"
	"k8s.io/klog"
)

func addThought(ctx context.Context, dc nykya.Config, opts AddOptions) error {
	klog.Infof("addNote %+v", opts)

	i := nykya.RenderInput{
		FrontMatter: nykya.FrontMatter{
			Kind:   opts.Kind,
			Posted: nykya.NewYAMLTime(opts.Timestamp),
		},
		Format: nykya.Markdown,
	}

	var err error
	if opts.Content != "" {
		i.Inline = opts.Content
	} else {
		i, err = openEditor(ctx, dc, i)
		if err != nil {
			return fmt.Errorf("openEditor: %w", err)
		}
	}

	outPath, err := localPath(dc, i.FrontMatter)
	if err != nil {
		return fmt.Errorf("local path: %w", err)
	}

	return saveItem(ctx, dc, i, filepath.Join(dc.In, outPath))
}