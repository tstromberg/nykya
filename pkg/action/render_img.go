package action

import (
	"fmt"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/tstromberg/daily/pkg/daily"
	"github.com/tstromberg/daily/pkg/fs"
	"k8s.io/klog"
)

// ThumbOpts are thumbnail soptions
type ThumbOpts struct {
	X       int
	Y       int
	Quality int
}

var defaultThumbOpts = []ThumbOpts{
	{X: 300, Y: 200, Quality: 70},
	{X: 800, Y: 600, Quality: 80},
	{X: 1920, Y: 1080, Quality: 85},
}

func renderJPEG(i *daily.Item, destRoot string) (*RenderedPost, error) {
	rp := &RenderedPost{
		Item:       i,
		Thumbnails: map[string]ThumbOpts{},
	}

	dest := filepath.Join(destRoot, i.Hier, filepath.Base(i.Source))
	err := fs.Copy(i.Source, dest)
	if err != nil {
		return rp, err
	}

	thumbDir := filepath.Join(destRoot, i.Hier, ".t")
	for _, t := range defaultThumbOpts {
		out, err := generateThumbnail(i.Source, thumbDir, t)
		if err != nil {
			return rp, err
		}
		rp.Thumbnails[out] = t
	}
	return rp, nil
}

func generateThumbnail(in string, thumbDir string, t ThumbOpts) (string, error) {
	thumbDest := filepath.Join(thumbDir, fmt.Sprintf("%dx%d@%d.jpg", t.X, t.Y, t.Quality))
	img, err := imgio.Open(in)
	if err != nil {
		return "", fmt.Errorf("open: %w", err)
	}

	thumb := transform.Resize(img, 800, 800, transform.Linear)
	klog.Infof("writing to output.jpg")

	if err := imgio.Save(thumbDest, thumb, imgio.JPEGEncoder(thumbQuality)); err != nil {
		return "", fmt.Errorf("save: %w", err)
	}
	return thumbDest, nil
}
