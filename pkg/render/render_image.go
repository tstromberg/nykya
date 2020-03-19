package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/otiai10/copy"
	"github.com/tstromberg/daily/pkg/daily"
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

func renderImage(i *daily.Item, destRoot string) (*RenderedPost, error) {
	rp := &RenderedPost{
		Item:   i,
		URL:    filepath.ToSlash(i.RelPath),
		Thumbs: map[string]ThumbMeta{},
	}

	fullDest := filepath.Join(destRoot, i.RelPath)
	err := copy.Copy(i.Path, fullDest)
	if err != nil {
		return rp, err
	}

	img, err := imgio.Open(i.Path)
	if err != nil {
		return rp, fmt.Errorf("imgio: %w", err)
	}
	ratio := float32(img.Bounds().Dx()) / float32(img.Bounds().Dy())
	klog.Infof("%s ratio (x=%d, y=%d): %2.f", i.Path, img.Bounds().Dx(), img.Bounds().Dy(), ratio)
	thumbDir := filepath.Join(destRoot, filepath.Dir(fullDest), ".t")
	if err := os.MkdirAll(thumbDir, 0600); err != nil {
		return rp, err
	}
	for name, t := range defaultThumbOpts {
		base := strings.Split(filepath.Base(i.Path), ".")[0]
		y := int(float32(t.X) / ratio)
		thumbDest := filepath.Join(thumbDir, fmt.Sprintf("%s_%dx%d@%d.jpg", base, t.X, y, t.Quality))
		klog.Infof("thumb %s (y=%d): %s", name, y, thumbDest)
		resized := transform.Resize(img, t.X, y, transform.Linear)

		// TODO: avoid doing work over again
		if err := imgio.Save(thumbDest, resized, imgio.JPEGEncoder(t.Quality)); err != nil {
			return rp, fmt.Errorf("save: %w", err)
		}
		rp.Thumbs[name] = ThumbMeta{
			X:    resized.Bounds().Dx(),
			Y:    resized.Bounds().Dy(),
			Path: thumbDest,
		}
	}
	return rp, nil
}
