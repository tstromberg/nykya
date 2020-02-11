package action

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"

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
	if p.Kind == "jpeg" {
		copyFile(p.Source, filepath.Join(dst, p.Hier))
	}
	return &RenderedPost{Metadata: p}, nil
}

func generateThumbnails(path string) error {
	img, err := imgio.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	thumb := transform.Resize(img, 800, 800, transform.Linear)
	klog.Infof("writing to output.jpg")

	if err := imgio.Save("output.jpg", thumb, imgio.JPEGEncoder(thumbQuality)); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}
