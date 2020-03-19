package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v1"
	"k8s.io/klog"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/tstromberg/daily/pkg/daily"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func fromYAML(path string) (*daily.Item, error) {
	klog.Infof("yaml: %s", path)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fm daily.FrontMatter
	err = yaml.Unmarshal(b, &fm)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	i := &daily.Item{
		FrontMatter: fm,
		Path:        path,
	}

	// TODO: Find a more elegant way to handle front-matter
	si := bytes.Index(b, []byte(daily.MarkdownSeparator))
	if si > 0 {
		i.Content = string(b[si+len(daily.MarkdownSeparator):])
	}

	klog.Infof("read: %+v", i)
	return i, nil
}

func fromJPEG(path string) (*daily.Item, error) {
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
	return &daily.Item{
		FrontMatter: daily.FrontMatter{
			Kind:   "jpeg",
			Posted: daily.NewYAMLTime(t),
		},
		Path: path,
	}, nil
}

func fromDirectory(path string, root string) ([]*daily.Item, error) {
	klog.Infof("Looking inside %s ...", path)
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*daily.Item
	for _, f := range fs {
		if f.IsDir() {
			dirItems, err := fromDirectory(filepath.Join(root, path, f.Name()), root)
			if err != nil {
				klog.Warningf("from dir %s: %v", f.Name(), err)
			}
			ps = append(ps, dirItems...)
			continue
		}

		klog.Infof("found %s", f.Name())
		fp := filepath.Join(path, f.Name())
		p, err := fromFile(fp)
		if err != nil {
			klog.Warningf("unable to parse %s: %v", fp, err)
			continue
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func fromFile(path string) (*daily.Item, error) {
	klog.Infof("parsing: %v", path)
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return fromJPEG(path)
	case ".yaml":
		return fromYAML(path)
	default:
		return nil, fmt.Errorf("unknown file type: %q", ext)
	}
}
