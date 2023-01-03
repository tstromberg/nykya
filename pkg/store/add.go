package store

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"

	"github.com/tstromberg/nykya/pkg/nykya"
	"github.com/tstromberg/wotd/pkg/wotd"
)

// AddOptions are options that can be passed to the add command
type AddOptions struct {
	// Title is the title of the post
	Title string
	// Content is the content: may be a string or filename
	Content string
	// Kind is the kind of content (thought, post, image)
	Kind string
	// Format is the format of the content (JPEG, HTML, etc)
	Format string
	// Timestamp is when the content was posted
	Timestamp time.Time
}

// Add an object from a path
func Add(ctx context.Context, dc nykya.Config, opts AddOptions) error {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}

	switch opts.Kind {
	case "thought":
		return addThought(ctx, dc, opts)
	case "post":
		return addPost(ctx, dc, opts)
	case "image":
		return addImage(ctx, dc, opts)
	default:
		return fmt.Errorf("object type not in 'thought' or 'post': %q", opts.Kind)
	}
}

// saveItem save an item and all dependencies to disk (sidecars, images)
func saveItem(ctx context.Context, dc nykya.Config, i nykya.RenderInput, path string) error {
	// Save a local copy of non-inlined content
	if i.Inline == "" && i.FrontMatter.Origin != "" {
		relPath, err := localCopy(dc, i.FrontMatter)
		if err != nil {
			return fmt.Errorf("local copy: %w", err)
		}
		i.ContentPath = relPath
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	klog.Infof("marshalling: %+v", i)

	fm, err := yaml.Marshal(i.FrontMatter)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	switch i.Format {
	case nykya.Markdown:
		return saveMarkdown(path, fm, i.Inline)
	case nykya.HTML:
		return saveHTML(path, fm, i.Inline)
	case nykya.JPEG:
		return saveJPEG(path, fm)
	default:
		return fmt.Errorf("unknown format: %s", i.Format)
	}
}

func wordOfTheDay() (string, error) {
	rs, err := wotd.All(context.Background())
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	sb.WriteString("\n# words of the day\n#\n")

	for _, r := range rs {
		sb.WriteString(fmt.Sprintf("# %s —  %s\n", r.Word, r.URL))
		for _, d := range r.Definitions {
			for _, ps := range d.Parts {
				sb.WriteString(fmt.Sprintf("#    %-9.9s: %s\n", ps.Kind, ps.Text))
			}
			break
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func saveMarkdown(path string, fm []byte, content string) error {
	klog.Infof("saving markdown to %s ...", path)
	b := bytes.NewBuffer(fm)

	if len(content) == 0 {
		wotd, err := wordOfTheDay()
		if err != nil {
			klog.Errorf("wotd failure: %v", err)
		}

		_, err = b.WriteString(wotd)
		if err != nil {
			klog.Errorf("wotd append failed: %v", err)
		}
	}

	_, err := b.WriteString(nykya.MarkdownSeparator)
	if err != nil {
		return err
	}

	_, err = b.WriteString(content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b.Bytes(), 0644)
}

func saveHTML(path string, fm []byte, content string) error {
	klog.Infof("saving HTML to %s ...", path)
	b := bytes.NewBuffer([]byte(nykya.HTMLBegin))

	_, err := b.Write(fm)
	if err != nil {
		return err
	}

	_, err = b.WriteString(nykya.HTMLSeparator)
	if err != nil {
		return err
	}

	_, err = b.WriteString(content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b.Bytes(), 0644)
}

func saveJPEG(path string, fm []byte) error {
	return ioutil.WriteFile(path+".yaml", fm, 0644)
}
