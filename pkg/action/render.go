package action

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tstromberg/daily/pkg/parse"
	"github.com/tstromberg/daily/pkg/tmpl"
	"k8s.io/klog"
)

var (
	thumbQuality = 85
)

// Render takes an input subdirectory of objects and generates static output within another directory
func Render(src string, dst string) ([]string, error) {
	ps, err := parse.Root(src)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	klog.Infof("rendering %q to %q", src, dst)

	idx := filepath.Join(dst, "index.html")
	f, err := os.Create(idx)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	rps := []*RenderedPost{}
	for _, p := range ps {
		rp, err := renderPost(p, dst)
		if err != nil {
			klog.Errorf("renderPost(%+v): %v", p, err)
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

func renderPost(p *parse.Post, dst string) (*RenderedPost, error) {
	klog.Infof("render %+v to %s", p, dst)
	var err error
	if p.Kind == "jpeg" {
		return renderJPEG(p, dst)
	}
	return &RenderedPost{Metadata: p}, err
}
