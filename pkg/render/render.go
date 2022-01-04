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

	"github.com/otiai10/copy"
	"github.com/tstromberg/nykya/pkg/nykya"

	"k8s.io/klog/v2"
)

var thumbQuality = 85

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

	Previous    *nykya.RenderInput
	PreviousURL string

	Next    *nykya.RenderInput
	NextURL string
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

	sort.Slice(items, func(i, j int) bool {
		ip := items[i].FrontMatter.Date.Time
		jp := items[j].FrontMatter.Date.Time
		return ip.Before(jp)
	})

	for x, i := range items {
		ri, err := renderItem(ctx, dc, items, x)
		if err != nil {
			klog.Errorf("annotate(%+v): %v", i, err)
			continue
		}

		rs = append(rs, ri)
	}

	// Render indexes

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

	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
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
	dst := filepath.Join(dc.Out, "nykya_index.html")
	return dst, siteTmpl("index", dc.Theme, dst, st)
}

func renderItem(ctx context.Context, dc nykya.Config, is []*nykya.RenderInput, idx int) (*RenderedItem, error) {
	i := is[idx]
	klog.Infof("render: %s: %+v", i.ContentPath, i.FrontMatter)

	var previous, next *nykya.RenderInput
	if idx > 0 {
		previous = is[idx-1]
	}

	if idx < len(is)-1 {
		next = is[idx+1]
	}

	switch i.FrontMatter.Kind {
	case "image":
		return renderImage(ctx, dc, i)
	case "post":
		return renderPost(ctx, dc, i, previous, next)
	default:
		return renderRaw(ctx, dc, i)
	}
}

func renderRaw(ctx context.Context, dc nykya.Config, i *nykya.RenderInput) (*RenderedItem, error) {
	ri, _, err := copyRawFile(ctx, dc, i)
	return ri, err
}

func copyRawFile(ctx context.Context, dc nykya.Config, i *nykya.RenderInput) (*RenderedItem, bool, error) {
	ri := &RenderedItem{
		Input:   i,
		URL:     filepath.ToSlash(i.ContentPath),
		OutPath: i.ContentPath,
	}

	fullSrc := filepath.Join(dc.In, i.ContentPath)
	fullDest := filepath.Join(dc.Out, i.ContentPath)

	sst, err := os.Stat(fullSrc)
	if err != nil {
		return nil, false, err
	}

	dst, err := os.Stat(fullDest)
	updated := false

	if err != nil {
		updated = true
		klog.Infof("updating %s: does not exist", fullDest)
	}

	if err == nil && sst.Size() != dst.Size() {
		updated = true
		klog.Infof("updating %s: size mismatch", fullDest)
	}

	if err == nil && sst.ModTime().After(dst.ModTime()) {
		klog.Infof("updating %s: source newer", fullDest)
		updated = true
	}

	if updated {
		klog.Infof("copying %s to %s ...", fullSrc, fullDest)
		err = copy.Copy(fullSrc, fullDest)
	}

	return ri, updated, err
}
