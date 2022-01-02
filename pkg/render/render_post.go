package render

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tstromberg/nykya/pkg/nykya"
	"github.com/yuin/goldmark"
	"k8s.io/klog/v2"
)

func htmlContent(in string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(in))
	if err != nil {
		return "", fmt.Errorf("goquery: %w", err)
	}

	content := ""
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		content = s.Find("main").Text()
		if content != "" {
			klog.V(1).Infof("found <main>: %s", content)
			return
		}
		content = s.Text()
		klog.V(1).Infof("found <body>: %s", content)
	})

	if content == "" {
		klog.Warningf("did not find body tag: %s", content)
	}
	return content, nil
}

func markdownContent(in string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(in), &buf); err != nil {
		return "", fmt.Errorf("goldmark: %w", err)
	}
	return buf.String(), nil
}

func renderPost(ctx context.Context, dc nykya.Config, i *nykya.RenderInput) (*RenderedItem, error) {
	ext := filepath.Ext(i.ContentPath)
	outPath := strings.Replace(i.ContentPath, ext, ".html", 1)

	ri := &RenderedItem{
		Title:   i.FrontMatter.Title,
		Input:   i,
		URL:     filepath.ToSlash(outPath),
		OutPath: outPath,
	}

	var err error
	var content string
	switch i.Format {
	case nykya.HTML:
		content, err = htmlContent(i.Inline)
	case nykya.Markdown:
		content, err = markdownContent(i.Inline)
	default:
		return ri, fmt.Errorf("unknown format: %q", i.Format)
	}

	if err != nil {
		return ri, err
	}

	ri.Content = template.HTML(content)

	klog.V(1).Infof("%s content: %s", ri.Title, content)
	return ri, siteTmpl("post", dc.Theme, filepath.Join(dc.Out, outPath), ri)
}
