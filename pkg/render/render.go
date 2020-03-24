package render

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tstromberg/nykya/pkg/nykya"

	"k8s.io/klog"
)

var (
	thumbQuality = 85
)

// Stream is basically the entire blog.
type Stream struct {
	Rendered    []*RenderedItem
	Title       string
	Subtitle    string
	Description string

	Timestamp time.Time
}

// RenderedItem is a post along with any dynamically generated data we found
type RenderedItem struct {
	RawItem *nykya.RawItem
	OutPath string
	URL     string
	Thumbs  map[string]ThumbMeta

	Content template.HTML
	Title   string
}

func indexesForRawItem(i *nykya.RawItem) []string {
	// TODO: make this more advanced
	base := strings.Split(filepath.ToSlash(i.RelPath), "/")[0]
	return []string{"/", base}
}

// Site generates static output to the site output directory
func Site(ctx context.Context, dc nykya.Config, items []*nykya.RawItem) ([]string, error) {
	rs := []*RenderedItem{}
	paths := []string{}

	for _, i := range items {
		ri, err := renderItem(ctx, dc, i)
		if err != nil {
			klog.Errorf("annotate(%+v): %v", i, err)
			continue
		}

		rs = append(rs, ri)
	}

	st := &Stream{
		Title:       dc.Title,
		Subtitle:    dc.Subtitle,
		Description: dc.Description,
		Timestamp:   time.Now(),
		Rendered:    rs,
	}

	path, err := siteIndex(ctx, dc, st)
	if err != nil {
		return []string{path}, fmt.Errorf("site index: %w", err)
	}
	paths = append(paths, path)
	return paths, nil
}

func siteTmpl(name string, themeRoot string, dst string, data interface{}) error {
	klog.Infof("Rendering %s to %s: %+v", name, dst, data)

	if err := os.MkdirAll(filepath.Dir(dst), 0600); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	paths := []string{
		filepath.Join(themeRoot, fmt.Sprintf("%s.tmpl", name)),
		filepath.Join(themeRoot, "style.tmpl"),
		filepath.Join(themeRoot, "footer.tmpl"),
		filepath.Join(themeRoot, "js.tmpl"),
		filepath.Join(themeRoot, "base.tmpl"),
	}

	t := template.Must(template.New(fmt.Sprintf("%s.tmpl", name)).ParseFiles(paths...))
	err = t.Execute(f, data)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	return nil

}

func siteIndex(ctx context.Context, dc nykya.Config, st *Stream) (string, error) {
	dst := filepath.Join(dc.Out, "index.html")
	return dst, siteTmpl("index", dc.Theme, dst, st)
}

func renderItem(ctx context.Context, dc nykya.Config, i *nykya.RawItem) (*RenderedItem, error) {
	klog.Infof("render %s %s: %+v", i.FrontMatter.Kind, i.RelPath, i)

	switch i.FrontMatter.Kind {
	case "image":
		return image(ctx, dc, i)
	case "post":
		return post(ctx, dc, i)
	default:
		return &RenderedItem{
			RawItem: i,
			Content: template.HTML(i.Content),
		}, nil
	}
}
