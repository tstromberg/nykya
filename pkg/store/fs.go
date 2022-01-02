package store

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"github.com/tstromberg/nykya/pkg/nykya"
	"k8s.io/klog/v2"
)

// localPath calculates the local path of an input, returning the relative path
func localPath(dc nykya.Config, fm nykya.FrontMatter) (string, error) {
	relSrc, err := filepath.Rel(dc.In, fm.Origin)

	if err == nil && !filepath.IsAbs(relSrc) && !strings.Contains(relSrc, "..") {
		klog.Infof("%s looks relative", fm.Origin)
		return relSrc, nil
	}

	relDir, err := calculateInputHierarchy(dc, fm)
	if err != nil {
		return "", fmt.Errorf("calculate hierarchy: %w", err)
	}

	newPath := filepath.Join(relDir, filepath.Base(fm.Origin))
	klog.Infof("new path for %s: %s", fm.Origin, newPath)
	return newPath, nil
}

// localCopy makes a local copy of an input, returning the relative path
func localCopy(dc nykya.Config, fm nykya.FrontMatter) (string, error) {
	if fm.Origin == "" {
		return "", nil
	}

	relDest, err := localPath(dc, fm)
	if err != nil {
		return "", fmt.Errorf("local path: %w", err)
	}

	fullDest := filepath.Join(dc.In, relDest)
	klog.Infof("copying %s -> %s", fm.Origin, fullDest)
	return relDest, copy.Copy(fm.Origin, fullDest)
}

func extForFormat(f string) string {
	switch f {
	case nykya.Markdown:
		return ".md"
	case nykya.HTML:
		return ".html"
	default:
		return "." + strings.ToLower(f)
	}
}

func formatForPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md":
		return nykya.Markdown
	case ".html":
		return nykya.HTML
	case ".jpeg", ".jpg":
		return nykya.JPEG
	default:
		return strings.Replace(ext, ".", "", 1)
	}
}

// calculateInputHierarchy calculates the relative destination directory for a file
func calculateInputHierarchy(dc nykya.Config, fm nykya.FrontMatter) (string, error) {
	tmpl := dc.Organization[fm.Kind]
	if tmpl == "" {
		tmpl = nykya.DefaultOrganization
	}

	t, err := template.New("orgtmpl").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parsing %q: %w", tmpl, err)
	}

	var b bytes.Buffer
	err = t.Execute(&b, fm)
	if err != nil {
		return "", fmt.Errorf("execute: %w", err)
	}

	return b.String(), nil
}
