package render

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/otiai10/copy"
	"github.com/tstromberg/paivalehti/pkg/paivalehti"
	"k8s.io/klog"
)

// ThumbOpts are thumbnail soptions
type ThumbOpts struct {
	X       int
	Y       int
	Quality int
}

// ThumbMeta describes a thumbnail
type ThumbMeta struct {
	X    int
	Y    int
	Path string
}

var defaultThumbOpts = map[string]ThumbOpts{
	"100w":  ThumbOpts{X: 100, Quality: 70},
	"200w":  ThumbOpts{X: 200, Quality: 70},
	"400w":  ThumbOpts{X: 400, Quality: 70},
	"800w":  ThumbOpts{X: 800, Quality: 80},
	"1920w": ThumbOpts{X: 1920, Quality: 85},
}

func image(ctx context.Context, dc paivalehti.Config, i *paivalehti.Item) (*renderedItem, error) {
	ri := &renderedItem{
		Item:    i,
		URL:     filepath.ToSlash(i.RelPath),
		OutPath: i.RelPath,
		Thumbs:  map[string]ThumbMeta{},
	}

	fullDest := filepath.Join(dc.Out, i.RelPath)
	err := copy.Copy(i.Path, fullDest)
	if err != nil {
		return ri, err
	}

	img, err := imgio.Open(i.Path)
	if err != nil {
		return ri, fmt.Errorf("imgio: %w", err)
	}
	ratio := float32(img.Bounds().Dx()) / float32(img.Bounds().Dy())
	klog.Infof("%s ratio (x=%d, y=%d): %2.f", i.Path, img.Bounds().Dx(), img.Bounds().Dy(), ratio)
	thumbDir := filepath.Join(filepath.Dir(fullDest), ".t")
	if err := os.MkdirAll(thumbDir, 0600); err != nil {
		return ri, err
	}
	for name, t := range defaultThumbOpts {
		base := strings.Split(filepath.Base(i.Path), ".")[0]
		y := int(float32(t.X) / ratio)
		thumbDest := filepath.Join(thumbDir, fmt.Sprintf("%s_%dx%d@%d.jpg", base, t.X, y, t.Quality))
		klog.Infof("thumb %s (y=%d): %s", name, y, thumbDest)
		resized := transform.Resize(img, t.X, y, transform.Linear)

		// TODO: avoid doing work over again
		if err := imgio.Save(thumbDest, resized, imgio.JPEGEncoder(t.Quality)); err != nil {
			return ri, fmt.Errorf("save: %w", err)
		}
		ri.Thumbs[name] = ThumbMeta{
			X:    resized.Bounds().Dx(),
			Y:    resized.Bounds().Dy(),
			Path: thumbDest,
		}
	}
	return ri, nil
}
