package render

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tstromberg/nykya/pkg/nykya"
	"github.com/yuin/goldmark"
	"k8s.io/klog/v2"
)

var (
	BareImgRe    = regexp.MustCompile(`\[img \"(.*?)\"\]`)
	CaptionImgRe = regexp.MustCompile(`\[img \"(.*?)\" \"(.*?)\"\]`)
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

	s := buf.String()
	s = BareImgRe.ReplaceAllString(s, `<figure><img src="$1"></figure>`)
	s = CaptionImgRe.ReplaceAllString(s, `<figure><img src="$1" alt="$2"><figcaption>$2</figcaption></figure>`)
	return s, nil
}

func urlTo(i *nykya.RenderInput) string {
	if i == nil {
		return ""
	}

	ext := filepath.Ext(i.ContentPath)
	url := filepath.ToSlash(strings.Replace(i.ContentPath, ext, ".html", 1))
	url = strings.Replace(url, "/post.html", ".html", 1)
	url = strings.Replace(url, "/index.html", ".html", 1)
	return url
}

func relPath(i *nykya.RenderInput, k *nykya.RenderInput) string {
	src := urlTo(i)
	dest := urlTo(k)
	rel, err := filepath.Rel(src, dest)
	if err != nil {
		return fmt.Sprintf("error.%v", err)
	}
	// lame
	return filepath.Base(rel)
}

func renderPost(ctx context.Context, dc nykya.Config, i *nykya.RenderInput, previous *nykya.RenderInput, next *nykya.RenderInput) (*RenderedItem, error) {
	ext := filepath.Ext(i.ContentPath)

	outPath := strings.Replace(i.ContentPath, ext, ".html", 1)
	outPath = strings.Replace(outPath, "/post.html", ".html", 1)
	outPath = strings.Replace(outPath, "/index.html", ".html", 1)

	ri := &RenderedItem{
		PageTitle:   i.FrontMatter.Title,
		Input:       i,
		SiteTitle:   dc.Title,
		URL:         urlTo(i),
		OutPath:     outPath,
		Next:        next,
		Previous:    previous,
		PreviousURL: relPath(i, previous),
		NextURL:     relPath(i, next),
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

	klog.V(1).Infof("%s content: %s", ri.PageTitle, content)
	return ri, siteTmpl("post", dc, filepath.Join(dc.Out, outPath), ri)
}
