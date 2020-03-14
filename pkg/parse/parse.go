package parse

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"k8s.io/klog"

	"github.com/tstromberg/daily/pkg/daily"
)

func Scan(ctx context.Context, root string) ([]*daily.Item, error) {
	klog.Infof("Scanning root %s ...", root)

	fs, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*daily.Item

	for _, f := range fs {
		klog.Infof("Scanning subdir %s", f.Name())
		ds, err := fromDirectory(filepath.Join(root, f.Name()), root)
		if err != nil {
			return nil, fmt.Errorf("parse date: %w", err)
		}
		ps = append(ps, ds...)
	}
	return ps, nil
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
		if p.Hier == "" {
			rel, err := filepath.Rel(filepath.Dir(path), root)
			if err != nil {
				return ps, fmt.Errorf("relpath: %w", err)
			}
			p.Hier = rel
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
