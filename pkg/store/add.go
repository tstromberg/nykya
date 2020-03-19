package store

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tstromberg/paivalehti/pkg/paivalehti"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

// AddOptions are options that can be passed to the add command
type AddOptions struct {
	Root string

	Title   string
	Content string
	Kind    string
	Source  string

	Timestamp time.Time
}

// Add an object from a path
func Add(ctx context.Context, dc paivalehti.Config, opts AddOptions) error {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	klog.Infof("adding %q to %s, ts=%s opts=%+v", opts.Content, opts.Root, ts, opts)

	switch opts.Kind {
	case "thought":
		return addThought(ctx, dc, opts)
	case "post":
		return addPost(ctx, dc, opts)
	default:
		return fmt.Errorf("unknown object type: %q", opts.Kind)
	}
}

// add is for adding thoughts
func addThought(ctx context.Context, dc paivalehti.Config, opts AddOptions) error {
	klog.Infof("addNote %+v", opts)

	words := strings.Split(strings.ToLower(opts.Content), " ")
	slug := strings.Join(words[0:3], "-")

	i := paivalehti.Item{
		FrontMatter: paivalehti.FrontMatter{
			Kind:   opts.Kind,
			Posted: paivalehti.NewYAMLTime(opts.Timestamp),
			Source: opts.Source,
		},
		Content: opts.Content,
		Format:  paivalehti.Markdown,
	}

	od, err := inDir(dc, i.FrontMatter)
	if err != nil {
		return fmt.Errorf("out dir(%+v): %w", i, err)
	}
	i.Path = filepath.Join(od, slug+".md")
	return saveItem(ctx, dc, i)
}

func extForFormat(f string) string {
	switch f {
	case paivalehti.Markdown:
		return ".md"
	case paivalehti.HTML:
		return ".html"
	default:
		return "." + strings.ToLower(f)
	}
}

func formatForPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md":
		return paivalehti.Markdown
	case ".html":
		return paivalehti.HTML
	default:
		return ext
	}
}

// addPost is for adding a post
func addPost(ctx context.Context, dc paivalehti.Config, opts AddOptions) error {
	// A post can be markdown, or HTML.
	// The file may, or may not exist.
	klog.Infof("addPost %+v", opts)

	if opts.Source == "" {
		return fmt.Errorf("no path specified")
	}

	fm := paivalehti.FrontMatter{
		Kind:   opts.Kind,
		Posted: paivalehti.NewYAMLTime(opts.Timestamp),
		Source: opts.Source,
	}

	path := opts.Source
	// Not right.. filepath.Rel?
	if strings.HasPrefix(dc.Out, path) {
		od, err := inDir(dc, fm)
		if err != nil {
			return fmt.Errorf("out dir: %w", err)
		}
		path = filepath.Join(od, filepath.Base(opts.Source))
	}

	i := paivalehti.Item{
		Content: opts.Content,
		Path:    path,
		Format:  formatForPath(opts.Source),
	}
	return saveItem(ctx, dc, i)
}

// saveItem saves an item to disk
func saveItem(ctx context.Context, dc paivalehti.Config, i paivalehti.Item) error {
	klog.Infof("marshalling: %+v", i)
	b, err := yaml.Marshal(i.FrontMatter)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	klog.Infof(string(b))

	dir := filepath.Dir(i.Path)
	if _, err := os.Stat(dir); err != nil {
		klog.Infof("Creating %s ...", dir)
		err := os.MkdirAll(dir, 0600)
		if err != nil {
			klog.Errorf("mkdir(%s) failed: %v", dir, err)
		}
	}

	fmt.Printf("Writing to %s ...", i.Path)
	f, err := os.Create(i.Path)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	switch i.Format {
	case paivalehti.Markdown:
		return saveMarkdown(f, b, i.Content)
	case paivalehti.HTML:
		return saveHTML(f, b, i.Content)
	default:
		return fmt.Errorf("unknown format: %s", i.Format)
	}
}

// inDir calculates the input directory for a file
func inDir(dc paivalehti.Config, fm paivalehti.FrontMatter) (string, error) {
	tmpl := dc.Organization[fm.Kind]
	if tmpl == "" {
		tmpl = paivalehti.DefaultOrganization
	}
	klog.Infof("inDir for %s: root=%q in=%q tmpl=%q", fm.Kind, dc.Root, dc.In, tmpl)

	t, err := template.New("orgtmpl").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parsing %q: %w", tmpl, err)
	}

	var b bytes.Buffer
	err = t.Execute(&b, fm)
	if err != nil {
		return "", fmt.Errorf("execute: %w", err)
	}

	return filepath.Join(dc.In, b.String()), nil
}

func saveMarkdown(w io.Writer, bs []byte, content string) error {
	w.Write(bs)
	io.WriteString(w, paivalehti.MarkdownSeparator)
	_, err := io.WriteString(w, content)
	return err
}

func saveHTML(w io.Writer, bs []byte, content string) error {
	io.WriteString(w, paivalehti.HTMLBegin)
	w.Write(bs)
	io.WriteString(w, paivalehti.HTMLSeparator)
	_, err := io.WriteString(w, content)
	return err
}
