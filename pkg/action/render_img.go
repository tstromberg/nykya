package action

import (
	"fmt"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/tstromberg/daily/pkg/fs"
	"github.com/tstromberg/daily/pkg/parse"
	"k8s.io/klog"
)

func renderJPEG(p *parse.Post, destRoot string) (*RenderedPost, error) {
	rp := &RenderedPost{
		Metadata: p,
	}

	dest := filepath.Join(destRoot, p.Hier, filepath.Base(p.Source))
	err := fs.Copy(p.Source, dest)
	if err != nil {
		return rp, err
	}

	if err := generateThumbnails(dest); err != nil {
		return rp, err
	}

	return rp, nil
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
