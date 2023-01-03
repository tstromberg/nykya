package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/djherbis/times.v1"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"

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

func fromMarkdown(path string) (*nykya.RenderInput, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	b = normalizeNewlines(b)

	i := &nykya.RenderInput{
		FrontMatter: nykya.FrontMatter{
			Date: nykya.NewYAMLTime(t.ModTime()),
		},
		Format: nykya.Markdown,
	}

	if i.FrontMatter.Kind == "" {
		i.FrontMatter.Kind = "post"
	}

	before, after, found := bytes.Cut(b, []byte(nykya.MarkdownSeparator))
	if !found {
		i.Inline = string(b)
		return i, nil
	}

	err = yaml.Unmarshal(before, &i.FrontMatter)
	if err != nil {
		return nil, fmt.Errorf("unmarshal of %s: %w", b, err)
	}
	i.Inline = string(after)

	return i, nil
}

func fromHTML(path string, relPath string) (*nykya.RenderInput, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	b = normalizeNewlines(b)

	i := &nykya.RenderInput{
		FrontMatter: nykya.FrontMatter{
			Date: nykya.NewYAMLTime(t.ModTime()),
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
			klog.Infof("%s: found html content: %s", path, i.Inline)
			i.Inline = string(b[si+len(nykya.HTMLSeparator):])
		}
	} else {
		return fromRawFile(path, relPath)
	}

	return i, nil
}

func fromJPEG(path string) (*nykya.RenderInput, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	i := &nykya.RenderInput{
		FrontMatter: nykya.FrontMatter{
			Kind: "image",
			Date: nykya.NewYAMLTime(t.ModTime()),
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
		i.FrontMatter.Date = nykya.NewYAMLTime(et)
	}

	ed, err := ex.Get(exif.ImageDescription)
	if err == nil {
		i.FrontMatter.Description = ed.String()
	}

	return i, nil
}

func fromRawFile(path string, relPath string) (*nykya.RenderInput, error) {
	t, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	i := &nykya.RenderInput{
		FrontMatter: nykya.FrontMatter{
			Kind: "raw",
			Date: nykya.NewYAMLTime(t.ModTime()),
		},
		ContentPath: relPath,
		Format:      nykya.Raw,
	}

	return i, nil
}

func fromDirectory(path string, root string) ([]*nykya.RenderInput, error) {
	klog.Infof("Looking inside %s ...", path)
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	var ps []*nykya.RenderInput

	for _, f := range fs {
		if strings.HasPrefix(f.Name(), ".") {
			klog.Warningf("skipping %s (hidden)", f.Name())
			continue
		}

		if f.IsDir() {
			dirRenderInputs, err := fromDirectory(filepath.Join(path, f.Name()), root)
			if err != nil {
				klog.Errorf("%s returned error: %v", f.Name(), err)
			}
			ps = append(ps, dirRenderInputs...)
			continue
		}

		klog.Infof("  found %s", f.Name())
		fp := filepath.Join(path, f.Name())

		rel, err := filepath.Rel(root, fp)
		if err != nil {
			return ps, fmt.Errorf("rel: %w", err)
		}

		i, err := fromFile(fp, rel)
		// Mostly harmless
		if err != nil {
			klog.Errorf("fromFile failed on %q: %v", fp, err)
			continue
		}

		if i == nil {
			klog.Infof("ignoring %s (no output", f.Name())
			continue
		}

		klog.V(1).Infof("  date=%s title=%s", i.FrontMatter.Date, i.FrontMatter.Title)

		i.ContentPath = rel
		ps = append(ps, i)
	}
	return ps, nil
}

func fromFile(path string, relPath string) (*nykya.RenderInput, error) {
	klog.V(1).Infof("parsing: %v", path)

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return fromJPEG(path)
	case ".yaml":
		return nil, nil
	case ".md":
		return fromMarkdown(path)
	case ".html":
		return fromHTML(path, relPath)
	case ".DS_Store", ".ds_store":
		return nil, nil
	default:
		return fromRawFile(path, relPath)
	}
}
