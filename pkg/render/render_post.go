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
	"k8s.io/klog"
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
			klog.Infof("found <main>: %s", content)
			return
		}
		content = s.Text()
		klog.Infof("found <body>: %s", content)
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

func post(ctx context.Context, dc nykya.Config, i *nykya.RawItem) (*RenderedItem, error) {
	ext := filepath.Ext(i.RelPath)
	outPath := strings.Replace(i.RelPath, ext, ".html", 1)

	ri := &RenderedItem{
		Title:   i.FrontMatter.Title,
		RawItem: i,
		URL:     filepath.ToSlash(outPath),
		OutPath: outPath,
	}

	var err error
	var content string
	switch i.Format {
	case nykya.HTML:
		content, err = htmlContent(i.Content)
	case nykya.Markdown:
		content, err = markdownContent(i.Content)
	default:
		return ri, fmt.Errorf("unknown format: %q", i.Format)
	}

	if err != nil {
		return ri, err
	}

	ri.Content = template.HTML(content)

	klog.Infof("%s content: %s", ri.Title, content)
	return ri, siteTmpl("post", dc.Theme, filepath.Join(dc.Out, outPath), ri)
}
