package tree

import (
	"fmt"
	"strings"
	"time"
	"path/filepath"
	"io/ioutil"
	
	"k8s.io/klog"
)


func Parse(root string) (*Stream, error) {
	klog.Infof("Parsing %s ...", root)
	s := &Stream{
		Timestamp: time.Now(),
	}


	fs, err := ioutil.ReadDir(root)
	if err != nil {
		return s, fmt.Errorf("readdir: %w", err)
	}
    for _, f := range fs {
		klog.Infof("found %s", f.Name())
		ps, err := parseDaily(filepath.Join(root, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("parse date: %w", err)
		}
		s.Posts = append(s.Posts, ps...)
	}
	return s, nil
}

func parseDaily(path string) ([]*Post, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*Post
    for _, f := range fs {
		klog.Infof("found %s", f.Name())
		p, err := parseFile(filepath.Join(path, f.Name()))
		if err != nil {
			klog.Warningf("unable to parse %s: %v", path, err)
			continue
		}
		ps = append(ps, p)
	}
	return ps,nil
}

func parseFile(path string) (*Post, error) {
	klog.Infof("parsing: %v", path)
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
			return parseJPEG(path)
	default:
		return nil, fmt.Errorf("unknown file type: %q", ext)
	}
}
