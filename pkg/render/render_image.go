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
	"github.com/tstromberg/nykya/pkg/nykya"
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
	X   int
	Y   int
	Src string
}

var defaultThumbOpts = map[string]ThumbOpts{
	"100w":  {X: 100, Quality: 70},
	"200w":  {X: 200, Quality: 70},
	"400w":  {X: 400, Quality: 70},
	"800w":  {X: 800, Quality: 80},
	"1920w": {X: 1920, Quality: 85},
}

func image(ctx context.Context, dc nykya.Config, i *nykya.RenderInput) (*RenderedItem, error) {
	ri := &RenderedItem{
		Input:   i,
		URL:     filepath.ToSlash(i.ContentPath),
		OutPath: i.ContentPath,
		Thumbs:  map[string]ThumbMeta{},
	}

	fullSrc := filepath.Join(dc.In, i.ContentPath)
	fullDest := filepath.Join(dc.Out, i.ContentPath)

	err := copy.Copy(fullSrc, fullDest)
	if err != nil {
		return ri, err
	}

	img, err := imgio.Open(fullSrc)
	if err != nil {
		return ri, fmt.Errorf("imgio: %w", err)
	}

	ratio := float32(img.Bounds().Dx()) / float32(img.Bounds().Dy())
	klog.Infof("%s ratio (x=%d, y=%d): %2.f", fullSrc, img.Bounds().Dx(), img.Bounds().Dy(), ratio)

	thumbDir := filepath.Join(filepath.Dir(i.ContentPath), ".t")

	if err := os.MkdirAll(filepath.Join(dc.Out, thumbDir), 0600); err != nil {
		return ri, err
	}

	for name, t := range defaultThumbOpts {
		base := strings.Split(filepath.Base(i.ContentPath), ".")[0]
		y := int(float32(t.X) / ratio)

		thumbName := fmt.Sprintf("%s_%dx%d@%d.jpg", base, t.X, y, t.Quality)
		thumbDest := filepath.Join(thumbDir, thumbName)
		fullThumbDest := filepath.Join(dc.Out, thumbDest)

		st, err := os.Stat(fullThumbDest)
		rimg := transform.Resize(img, t.X, y, transform.Linear)

		ri.Thumbs[name] = ThumbMeta{
			X:   rimg.Bounds().Dx(),
			Y:   rimg.Bounds().Dy(),
			Src: thumbDest,
		}

		if err == nil && st.Size() > int64(128) {
			klog.Infof("%s exists", fullThumbDest)
			continue
		}

		if err := imgio.Save(fullThumbDest, rimg, imgio.JPEGEncoder(t.Quality)); err != nil {
			return ri, fmt.Errorf("save: %w", err)
		}
	}

	klog.Infof("image: %+v", ri)
	return ri, nil
}
