package action

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"

	"github.com/tstromberg/daily/pkg/tree"
	"github.com/tstromberg/daily/pkg/tmpl"
	"k8s.io/klog"
)

var (
	thumbQuality = 85
)

func Render(src string, dst string) ([]string, error) {
	st, err := tree.Parse(src)
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
	return []string{idx}, tmpl.Index.Execute(f, st)
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
