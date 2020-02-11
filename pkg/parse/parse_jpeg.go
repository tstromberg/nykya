package parse

import (
	"fmt"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"k8s.io/klog"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func fromJPEG(path string) (*Post, error) {
	klog.Infof("jpeg: %s", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	t, err := x.DateTime()
	if err != nil {
		klog.Errorf("datetime(%s): %v", path, err)
	}
	return &Post{
		Kind:        "jpeg",
		Source:      path,
		Created:     t,
		Description: "Just another day in paradise",
	}, nil
}
