package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/djherbis/times.v1"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/tstromberg/paivalehti/pkg/paivalehti"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func fromYAML(path string) (*paivalehti.Item, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	i := &paivalehti.Item{
		FrontMatter: paivalehti.FrontMatter{
			Posted: paivalehti.NewYAMLTime(t.ModTime()),
		},
	}

	err = yaml.Unmarshal(b, &i.FrontMatter)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	// TODO: Find a more elegant way to handle front-matter
	si := bytes.Index(b, []byte(paivalehti.MarkdownSeparator))
	if si > 0 {
		i.Content = string(b[si+len(paivalehti.MarkdownSeparator):])
	}
	return i, nil
}

func fromHTML(path string) (*paivalehti.Item, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	i := &paivalehti.Item{
		FrontMatter: paivalehti.FrontMatter{
			Posted: paivalehti.NewYAMLTime(t.ModTime()),
		},
	}

	header := b[0:len(paivalehti.HTMLBegin)]
	klog.V(1).Infof("%s header: %q vs %q", path, string(header), paivalehti.HTMLBegin)
	if bytes.Equal(header, []byte(paivalehti.HTMLBegin)) {
		si := bytes.Index(b, []byte(paivalehti.HTMLSeparator))
		if si > 0 {
			fb := b[len(header):si]
			klog.V(1).Infof("frontmatter bytes: %s", b)
			err = yaml.Unmarshal(fb, &i.FrontMatter)
			if err != nil {
				return nil, fmt.Errorf("unmarshal: %w", err)
			}
			i.Content = string(b[si+len(paivalehti.HTMLSeparator):])
		}
	}

	return i, nil
}

func fromJPEG(path string) (*paivalehti.Item, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	i := &paivalehti.Item{
		FrontMatter: paivalehti.FrontMatter{
			Kind:   "image",
			Posted: paivalehti.NewYAMLTime(t.ModTime()),
		},
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	ex, err := exif.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	et, err := ex.DateTime()
	if err != nil {
		klog.Errorf("datetime(%s): %v", path, err)
	} else {
		i.FrontMatter.Posted = paivalehti.NewYAMLTime(et)
	}

	ed, err := ex.Get(exif.ImageDescription)
	if err == nil {
		i.FrontMatter.Description = ed.String()
	}

	return i, nil
}

func fromDirectory(path string, root string) ([]*paivalehti.Item, error) {
	klog.V(2).Infof("Looking inside %s ...", path)
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*paivalehti.Item
	for _, f := range fs {
		if f.IsDir() {
			dirItems, err := fromDirectory(filepath.Join(root, path, f.Name()), root)
			if err != nil {
				klog.Warningf("from dir %s: %v", f.Name(), err)
			}
			ps = append(ps, dirItems...)
			continue
		}

		klog.V(1).Infof("found %s", f.Name())
		fp := filepath.Join(path, f.Name())

		rel, err := filepath.Rel(root, fp)
		if err != nil {
			return ps, fmt.Errorf("rel: %w", err)
		}

		i, err := fromFile(fp)
		if i.FrontMatter.Kind == "" {
			klog.Errorf("%s has no kind: %+v", fp, i)
			continue
		}
		i.Path = fp
		i.RelPath = rel

		if err != nil {
			klog.Warningf("unable to parse %s: %v", fp, err)
			continue
		}
		klog.Infof("%s == %s (%s)", rel, i.FrontMatter.Kind, i.FrontMatter.Title)
		ps = append(ps, i)
	}
	return ps, nil
}

func fromFile(path string) (*paivalehti.Item, error) {
	klog.V(1).Infof("parsing: %v", path)
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return fromJPEG(path)
	case ".yaml", ".md":
		return fromYAML(path)
	case ".html":
		return fromHTML(path)
	default:
		return nil, fmt.Errorf("unknown file type: %q", ext)
	}
}
