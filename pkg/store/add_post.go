package store

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tstromberg/nykya/pkg/nykya"
)

// addPost is for adding a post
func addPost(ctx context.Context, dc nykya.Config, opts AddOptions) error {
	i := nykya.RenderInput{
		FrontMatter: nykya.FrontMatter{
			Kind: opts.Kind,
			Date: nykya.NewYAMLTime(opts.Timestamp),
		},
		Format: opts.Format,
	}
	var err error
	var outPath string

	if opts.Content == "" {
		if i.Format == "" {
			i.Format = nykya.Markdown
		}

		i, err = openEditor(ctx, dc, i)
		if err != nil {
			return fmt.Errorf("openEditor: %w", err)
		}

		relDir, err := calculateInputHierarchy(dc, i.FrontMatter)
		if err != nil {
			return fmt.Errorf("calculate hierarchy: %w", err)
		}

		outPath = filepath.Join(relDir, fmt.Sprintf("%s.md", slugify(i.FrontMatter.Title)))
	} else {
		i.FrontMatter.Origin = opts.Content
		outPath, err = localPath(dc, i.FrontMatter)
		if err != nil {
			return fmt.Errorf("local path: %w", err)
		}
	}

	return saveItem(ctx, dc, i, filepath.Join(dc.In, outPath))
}
