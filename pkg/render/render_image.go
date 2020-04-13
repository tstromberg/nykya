package render

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/otiai10/copy"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/tstromberg/nykya/pkg/nykya"
	"k8s.io/klog"
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
	ri := &RenderedItem{
		Input:   i,
		URL:     filepath.ToSlash(i.ContentPath),
		OutPath: i.ContentPath,
		Thumbs:  map[string]ThumbMeta{},
	}

	fullSrc := filepath.Join(dc.In, i.ContentPath)
	fullDest := filepath.Join(dc.Out, i.ContentPath)

	sst, err := os.Stat(fullSrc)
	if err != nil {
		return nil, err
	}

	dst, err := os.Stat(fullDest)
	updated := false

	if err != nil {
		updated = true
		klog.Infof("updating %s: does not exist", fullDest)
	}

	if err == nil && sst.Size() != dst.Size() {
		updated = true
		klog.Infof("updating %s: size mismatch", fullDest)
	}

	if err == nil && sst.ModTime().After(dst.ModTime()) {
		klog.Infof("updating %s: source newer", fullDest)
		updated = true
	}

	if updated {
		err := copy.Copy(fullSrc, fullDest)
		if err != nil {
			return ri, err
		}
	}

	thumbDir := filepath.Join(filepath.Dir(i.ContentPath), ".t")
	if err := os.MkdirAll(filepath.Join(dc.Out, thumbDir), 0600); err != nil {
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
			klog.Warningf("unable to read thumb: %v", err)
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

	if t.X == 0 {
		scale := i.Bounds().Dy() / t.Y
		x = int(i.Bounds().Dx() / scale)
	}

	if t.Y == 0 {
		scale := i.Bounds().Dx() / t.X
		y = int(i.Bounds().Dy() / scale)
	}

	rimg := transform.Resize(i, x, y, transform.Lanczos)
	if err := imgio.Save(path, rimg, imgio.JPEGEncoder(t.Quality)); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &ThumbMeta{X: rimg.Bounds().Dx(), Y: rimg.Bounds().Dy()}, nil
}

func readThumb(path string) (*ThumbMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	ex, err := exif.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	gx, err := ex.Get(exif.ImageWidth)
	if err != nil {
		return nil, fmt.Errorf("imgwidth: %w", err)
	}
	x, err := strconv.Atoi(gx.String())
	if err != nil {
		return nil, err
	}

	gy, err := ex.Get(exif.ImageLength)
	if err != nil {
		return nil, fmt.Errorf("imglen: %w", err)
	}

	y, err := strconv.Atoi(gy.String())
	if err != nil {
		return nil, err
	}

	return &ThumbMeta{X: x, Y: y}, nil
}
