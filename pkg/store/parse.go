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
	"github.com/tstromberg/nykya/pkg/nykya"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func normalizeNewlines(bs []byte) []byte {
	bs = bytes.Replace(bs, []byte{13, 10}, []byte{10}, -1)
	return bytes.Replace(bs, []byte{13}, []byte{10}, -1)
}

func fromMarkdown(path string) (*nykya.RawItem, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	b = normalizeNewlines(b)

	i := &nykya.RawItem{
		FrontMatter: nykya.FrontMatter{
			Posted: nykya.NewYAMLTime(t.ModTime()),
		},
		Format: nykya.Markdown,
	}

	err = yaml.Unmarshal(b, &i.FrontMatter)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	si := bytes.Index(b, []byte(nykya.MarkdownSeparator))
	if si > 0 {
		i.Content = string(b[si+len(nykya.MarkdownSeparator):])
		klog.Infof("%s: found markdown content: %s", path, i.Content)
	} else {
		klog.Warningf("%s: did not find markdown content (si=%d)", path, si)
	}
	return i, nil
}

func fromHTML(path string) (*nykya.RawItem, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	b = normalizeNewlines(b)

	i := &nykya.RawItem{
		FrontMatter: nykya.FrontMatter{
			Posted: nykya.NewYAMLTime(t.ModTime()),
		},
		Format: nykya.HTML,
	}

	header := b[0:len(nykya.HTMLBegin)]
	klog.V(1).Infof("%s header: %q vs %q", path, string(header), nykya.HTMLBegin)
	if bytes.Equal(header, []byte(nykya.HTMLBegin)) {
		si := bytes.Index(b, []byte(nykya.HTMLSeparator))
		if si > 0 {
			fb := b[len(header):si]
			klog.V(1).Infof("frontmatter bytes: %s", b)
			err = yaml.Unmarshal(fb, &i.FrontMatter)
			if err != nil {
				return nil, fmt.Errorf("unmarshal: %w", err)
			}
			klog.Infof("%s: found html content: %s", path, i.Content)
			i.Content = string(b[si+len(nykya.HTMLSeparator):])
		}
	}

	return i, nil
}

func fromJPEG(path string) (*nykya.RawItem, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	i := &nykya.RawItem{
		FrontMatter: nykya.FrontMatter{
			Kind:   "image",
			Posted: nykya.NewYAMLTime(t.ModTime()),
		},
		Format: nykya.JPEG,
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
		i.FrontMatter.Posted = nykya.NewYAMLTime(et)
	}

	ed, err := ex.Get(exif.ImageDescription)
	if err == nil {
		i.FrontMatter.Description = ed.String()
	}

	return i, nil
}

func fromDirectory(path string, root string) ([]*nykya.RawItem, error) {
	klog.V(2).Infof("Looking inside %s ...", path)
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*nykya.RawItem
	for _, f := range fs {
		if f.IsDir() {
			dirRawItems, err := fromDirectory(filepath.Join(root, path, f.Name()), root)
			if err != nil {
				klog.Warningf("from dir %s: %v", f.Name(), err)
			}
			ps = append(ps, dirRawItems...)
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

func fromFile(path string) (*nykya.RawItem, error) {
	klog.V(1).Infof("parsing: %v", path)
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return fromJPEG(path)
	case ".yaml", ".md":
		return fromMarkdown(path)
	case ".html":
		return fromHTML(path)
	default:
		return nil, fmt.Errorf("unknown file type: %q", ext)
	}
}
