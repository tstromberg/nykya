package action

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tstromberg/daily/pkg/daily"
	"github.com/tstromberg/daily/pkg/parse"
	"github.com/tstromberg/daily/pkg/tmpl"

	"k8s.io/klog"
)

var (
	thumbQuality = 85
)

// Render takes an input subdirectory of objects and generates static output within another directory
func Render(ctx context.Context, dc daily.Config) ([]string, error) {
	items, err := parse.Scan(ctx, dc.In)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	idx := filepath.Join(dc.Out, "index.html")
	f, err := os.Create(idx)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	rps := []*RenderedPost{}
	for _, i := range items {
		rp, err := renderItem(ctx, i, dc.Out)
		if err != nil {
			klog.Errorf("renderPost(%+v): %v", i, err)
			continue
		}
		rps = append(rps, rp)
	}

	st := &Stream{
		Title:     "boring blog",
		Timestamp: time.Now(),
		Posts:     rps,
	}
	return []string{idx}, tmpl.Index.Execute(f, st)
}

func renderItem(ctx context.Context, i *daily.Item, dst string) (*RenderedPost, error) {
	klog.Infof("render %+v to %s", i, dst)
	var err error
	if i.Kind == "jpeg" {
		return renderJPEG(i, dst)
	}
	return &RenderedPost{Item: i}, err
}
