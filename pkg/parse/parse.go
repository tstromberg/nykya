package parse

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"k8s.io/klog"
)

func Root(root string) ([]*Post, error) {
	klog.Infof("Parsing %s ...", root)

	fs, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*Post

	for _, f := range fs {
		klog.Infof("found %s", f.Name())
		ds, err := fromDirectory(filepath.Join(root, f.Name()), root)
		if err != nil {
			return nil, fmt.Errorf("parse date: %w", err)
		}
		ps = append(ps, ds...)
	}
	return ps, nil
}

func fromDirectory(path string, root string) ([]*Post, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*Post
	for _, f := range fs {
		klog.Infof("found %s", f.Name())
		fp := filepath.Join(path, f.Name())
		p, err := fromFile(fp)
		if err != nil {
			klog.Warningf("unable to parse %s: %v", fp, err)
			continue
		}
		if p.Hierarchy == "" {
			rel, err := filepath.Rel(filepath.Dir(path), root)
			if err != nil {
				return ps, fmt.Errorf("relpath: %w", err)
			}
			p.Hierarchy = rel
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func fromFile(path string) (*Post, error) {
	klog.Infof("parsing: %v", path)
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return fromJPEG(path)
	default:
		return nil, fmt.Errorf("unknown file type: %q", ext)
	}
}
