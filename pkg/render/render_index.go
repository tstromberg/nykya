package render

import (
	"context"
	"path/filepath"
	"sort"
	"time"

	"github.com/tstromberg/nykya/pkg/nykya"
)

// Stream is basically the entire blog.
type Stream struct {
	Rendered  []*RenderedItem
	SiteTitle string
	PageTitle string
	OutPath   string
	Timestamp time.Time
}

func renderIndex(ctx context.Context, dc nykya.Config, rs []*RenderedItem, outPath string) error {
	// Newest first
	sort.Slice(rs, func(i, j int) bool {
		ip := rs[i].Input.FrontMatter.Date.Time
		jp := rs[j].Input.FrontMatter.Date.Time
		return ip.After(jp)
	})

	st := &Stream{
		SiteTitle: dc.Title,
		PageTitle: filepath.Base(filepath.Dir(outPath)),
		Timestamp: time.Now(),
		Rendered:  rs,
		OutPath:   outPath,
	}

	return siteTmpl("index", dc.Theme, filepath.Join(dc.Out, outPath), st)
}
