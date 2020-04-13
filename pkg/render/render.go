package render

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
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
	Input   *nykya.RenderInput
	OutPath string
	URL     string
	Thumbs  map[string]ThumbMeta

	Content template.HTML
	Title   string
}

func indexesForRenderInput(i *nykya.RenderInput) []string {
	// TODO: make this more advanced
	base := strings.Split(filepath.ToSlash(i.ContentPath), "/")[0]
	return []string{"/", base}
}

// Site generates static output to the site output directory
func Site(ctx context.Context, dc nykya.Config, items []*nykya.RenderInput) ([]string, error) {
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

	sort.Slice(rs, func(i, j int) bool {
		ip := rs[i].Input.FrontMatter.Posted.Time
		jp := rs[j].Input.FrontMatter.Posted.Time
		return ip.After(jp)
	})

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
	klog.V(1).Infof("Rendering %s to %s: %+v", name, dst, data)

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

func renderItem(ctx context.Context, dc nykya.Config, i *nykya.RenderInput) (*RenderedItem, error) {
	klog.Infof("render: %s (%s)", i.ContentPath, i.FrontMatter.Kind)
	switch i.FrontMatter.Kind {
	case "image":
		return renderImage(ctx, dc, i)
	case "post":
		return renderPost(ctx, dc, i)
	default:
		return &RenderedItem{
			Input:   i,
			Content: template.HTML(i.Inline),
		}, nil
	}
}
