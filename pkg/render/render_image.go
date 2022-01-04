package render

import (
	"context"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/tstromberg/nykya/pkg/nykya"
	"k8s.io/klog/v2"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

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
	"133t": {Y: 133, Quality: 85},
}

func renderImage(ctx context.Context, dc nykya.Config, i *nykya.RenderInput) (*RenderedItem, error) {
	// renderRaw takes care of copying the original file over
	ri, updated, err := copyRawFile(ctx, dc, i)
	if err != nil {
		return ri, err
	}
	ri.Thumbs = map[string]ThumbMeta{}

	fullDest := filepath.Join(dc.Out, i.ContentPath)

	thumbDir := filepath.Join(filepath.Dir(i.ContentPath), ".t")
	if err := os.MkdirAll(filepath.Join(dc.Out, thumbDir), 0o700); err != nil {
		return ri, err
	}

	base := strings.Split(filepath.Base(i.ContentPath), ".")[0]
	var img image.Image

	for name, t := range defaultThumbOpts {
		thumbName := fmt.Sprintf("%s@%s.jpg", base, name)
		thumbDest := filepath.Join(thumbDir, thumbName)
		fullThumbDest := filepath.Join(dc.Out, thumbDest)

		st, err := os.Stat(fullThumbDest)
		if err == nil && st.Size() > int64(128) && !updated {
			klog.Infof("%s exists (%d bytes)", fullThumbDest, st.Size())
			rt, err := readThumb(fullThumbDest)
			if err == nil {
				rt.Src = thumbDest
				ri.Thumbs[name] = *rt
				continue
			}
			klog.Warningf("unable to read %s: %v", fullThumbDest, err)
		}

		if img == nil {
			klog.Infof("opening %s ...", fullDest)
			img, err = imgio.Open(fullDest)
			if err != nil {
				return nil, err
			}
		}

		ct, err := createThumb(img, fullThumbDest, t)
		if err != nil {
			return nil, fmt.Errorf("create thumb: %w", err)
		}

		ct.Src = filepath.ToSlash(thumbDest)
		ri.Thumbs[name] = *ct
	}

	klog.Infof("image: %+v", ri)
	return ri, nil
}

func createThumb(i image.Image, path string, t ThumbOpts) (*ThumbMeta, error) {
	klog.Infof("creating thumb: %s", path)
	x := t.X
	y := t.Y

	if t.X == 0 && t.Y == 0 {
		return nil, fmt.Errorf("both dimensions cannot be zero: %+v", t)
	}

	if t.X == 0 {
		scale := math.Max(float64(i.Bounds().Dy()/t.Y), 1)
		x = int(float64(i.Bounds().Dx()) / scale)
	}

	if t.Y == 0 {
		scale := math.Max(float64(i.Bounds().Dx()/t.X), 1)
		y = int(float64(i.Bounds().Dy()) / scale)
	}

	klog.Infof("%+v resize result: x=%d y=%d", t, x, y)
	rimg := transform.Resize(i, x, y, transform.Lanczos)
	if err := imgio.Save(path, rimg, imgio.JPEGEncoder(t.Quality)); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	tm := &ThumbMeta{X: rimg.Bounds().Dx(), Y: rimg.Bounds().Dy()}
	return tm, nil
}

func readThumb(path string) (*ThumbMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	im, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &ThumbMeta{X: im.Width, Y: im.Height}, nil
}
