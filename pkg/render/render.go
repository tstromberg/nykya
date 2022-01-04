package render

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/otiai10/copy"
	"github.com/tstromberg/nykya/pkg/nykya"

	"k8s.io/klog/v2"
)

var thumbQuality = 85

// RenderedItem is a post along with any dynamically generated data we found
type RenderedItem struct {
	Input   *nykya.RenderInput
	OutPath string
	URL     string
	Thumbs  map[string]ThumbMeta

	Content   template.HTML
	SiteTitle string
	PageTitle string

	Previous    *nykya.RenderInput
	PreviousURL string

	Next    *nykya.RenderInput
	NextURL string
}

func indexesForRender(i *RenderedItem) []string {
	return []string{filepath.Dir(i.OutPath)}
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

	byIndex := map[string][]*RenderedItem{}

	for x, i := range items {
		if i.FrontMatter.Draft && !dc.IncludeDrafts {
			klog.Infof("Ignoring draft: %s", i.FrontMatter.Title)
			continue
		}

		ri, err := renderItem(ctx, dc, items, x)
		if err != nil {
			klog.Errorf("annotate(%+v): %v", i, err)
			continue
		}
		paths = append(paths, ri.OutPath)

		if ri.PageTitle != "" {
			for _, i := range indexesForRender(ri) {
				if byIndex[i] == nil {
					byIndex[i] = []*RenderedItem{}
				}
				klog.Infof("Adding %q to index %q", ri.PageTitle, i)
				byIndex[i] = append(byIndex[i], ri)
			}
		}

		rs = append(rs, ri)
	}

	for idx, ris := range byIndex {
		outPath := filepath.Join(idx, "index.html")
		if err := renderIndex(ctx, dc, ris, outPath); err != nil {
			return paths, fmt.Errorf("render index: %w", err)
		}
		paths = append(paths, outPath)
	}

	return paths, nil
}

func tmplRelPath(root string, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		klog.Errorf("unable to calculate relpath of root=%s path=%s: %v", root, path, err)
		return path
	}

	// If we're in the same directory, don't go up
	r = strings.TrimPrefix(r, "../")
	klog.Infof("relpath of root=%s path=%s: %s", root, path, r)
	return r
}

func siteTmpl(name string, themeRoot string, dst string, data interface{}) error {
	klog.Infof("Rendering %s to %s ...", name, dst)

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

	fm := template.FuncMap{
		"RelPath": tmplRelPath,
	}

	t := template.Must(template.New(fmt.Sprintf("%s.tmpl", name)).Funcs(fm).ParseFiles(paths...))
	err = t.Execute(f, data)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	return nil
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
