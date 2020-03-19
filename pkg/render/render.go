package render

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tstromberg/paivalehti/pkg/paivalehti"
	"github.com/tstromberg/paivalehti/pkg/tmpl"

	"k8s.io/klog"
)

var (
	thumbQuality = 85
)

// Stream is basically the entire blog.
type Stream struct {
	Posts       []*RenderedPost
	Title       string
	Subtitle    string
	Description string

	Timestamp time.Time
}

// RenderedPost is a post along with any dynamically generated data we found
type RenderedPost struct {
	Item   *paivalehti.Item
	URL    string
	Thumbs map[string]ThumbMeta
}

// Site generates static output to the site output directory
func Site(ctx context.Context, dc paivalehti.Config, items []*paivalehti.Item) ([]string, error) {
	klog.Infof("Rendering site to %s", dc.Out)
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
		Title:       dc.Title,
		Subtitle:    dc.Subtitle,
		Description: dc.Description,
		Timestamp:   time.Now(),
		Posts:       rps,
	}
	return []string{idx}, tmpl.Index.Execute(f, st)
}

func renderItem(ctx context.Context, i *paivalehti.Item, dst string) (*RenderedPost, error) {
	klog.Infof("render %s %s", i.FrontMatter.Kind, i.RelPath)
	var err error
	if i.FrontMatter.Kind == "image" {
		return renderImage(i, dst)
	}
	return &RenderedPost{
		Item: i,
		URL:  filepath.ToSlash(i.RelPath),
	}, err
}
