package render

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tstromberg/paivalehti/pkg/paivalehti"

	"k8s.io/klog"
)

var (
	thumbQuality = 85
)

// Stream is basically the entire blog.
type Stream struct {
	Annotated   []*annotatedItem
	Title       string
	Subtitle    string
	Description string

	Timestamp time.Time
}

// annotatedItem is a post along with any dynamically generated data we found
type annotatedItem struct {
	Item   *paivalehti.Item
	URL    string
	Thumbs map[string]ThumbMeta
}

func indexesForItem(i *paivalehti.Item) []string {
	// TODO: make this more advanced
	base := strings.Split(filepath.ToSlash(i.RelPath), "/")[0]
	return []string{"/", base}
}

// Site generates static output to the site output directory
func Site(ctx context.Context, dc paivalehti.Config, items []*paivalehti.Item) ([]string, error) {
	ais := []*annotatedItem{}
	paths := []string{}

	for _, i := range items {
		ai, err := annotate(ctx, i, dc.Out)
		if err != nil {
			klog.Errorf("annotate(%+v): %v", i, err)
			continue
		}

		ais = append(ais, ai)
	}

	st := &Stream{
		Title:       dc.Title,
		Subtitle:    dc.Subtitle,
		Description: dc.Description,
		Timestamp:   time.Now(),
		Annotated:   ais,
	}

	path, err := siteIndex(ctx, dc, st)
	if err != nil {
		return []string{path}, fmt.Errorf("site index: %w", err)
	}
	paths = append(paths, path)
	return paths, nil
}

func templatePaths(dc paivalehti.Config, name string) []string {
	return []string{
		filepath.Join(dc.Theme, fmt.Sprintf("%s.tmpl", name)),
		filepath.Join(dc.Theme, "style.tmpl"),
		filepath.Join(dc.Theme, "base.tmpl"),
	}
}

func siteIndex(ctx context.Context, dc paivalehti.Config, st *Stream) (string, error) {
	klog.Infof("Rendering site to %s", dc.Out)
	idx := filepath.Join(dc.Out, "index.html")
	f, err := os.Create(idx)
	if err != nil {
		return idx, fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	paths := templatePaths(dc, "index")
	klog.Infof("index files: %v", paths)
	t := template.Must(template.New("index.tmpl").ParseFiles(paths...))
	err = t.Execute(f, st)
	if err != nil {
		return idx, fmt.Errorf("execute: %w", err)
	}
	return idx, nil
}

func annotate(ctx context.Context, i *paivalehti.Item, dst string) (*annotatedItem, error) {
	klog.Infof("render %s %s", i.FrontMatter.Kind, i.RelPath)
	var err error
	if i.FrontMatter.Kind == "image" {
		return renderImage(i, dst)
	}
	return &annotatedItem{
		Item: i,
		URL:  filepath.ToSlash(i.RelPath),
	}, err
}
