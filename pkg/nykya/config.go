package nykya

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v1"
	"k8s.io/klog"
)

// DefaultOrganization shows where to put files if organization is unset
const DefaultOrganization = `{{ .Kind }}s/{{ .Posted.Format "2006-01-02" }}`

// Config is site configuration
type Config struct {
	Root string

	Title       string
	Subtitle    string
	Description string

	In    string
	Out   string
	Theme string

	Organization map[string]string
}

func findTheme(root string, theme string) string {
	pwd, err := os.Getwd()
	if err != nil {
		klog.Errorf("unable to getwd: %v", err)
		pwd = "."
	}
	if theme == "" {
		theme = "basic"
	}

	try := []string{
		theme,
		filepath.Join(root, "theme"),
		filepath.Join(root, theme),
		filepath.Join(pwd, theme),
		filepath.Join(root, "themes", theme),
		filepath.Join(pwd, "themes", theme),
		filepath.Join(root, "..", "themes", theme),
		filepath.Join(pwd, "..", "themes", theme),
		filepath.Join(root, "..", "..", "themes", theme),
		filepath.Join(pwd, "..", "..", "themes", theme),
	}

	for _, path := range try {
		_, err := os.Stat(filepath.Join(path, "base.tmpl"))
		if err == nil {
			klog.Infof("found theme: %s", path)
			return path
		}
		klog.Infof("tried %s", path)
	}
	return ""
}

// ConfigFromRoot returns the sites configuration
func ConfigFromRoot(root string) (Config, error) {
	c := Config{
		Root:        root,
		In:          filepath.Join(root, "in"),
		Out:         filepath.Join(root, "out"),
		Title:       "Example Title",
		Subtitle:    "Example Subtitle",
		Description: "Example Description",
	}

	root, err := filepath.Abs(root)
	if err != nil {
		return c, fmt.Errorf("abs: %w", err)
	}

	cp := filepath.Join(root, "nykya.yaml")
	if _, err := os.Stat(cp); err != nil {
		c.Theme = findTheme(root, c.Theme)
		klog.Infof("%s not found, returning demo site configuration", cp)
		return c, nil
	}

	b, err := ioutil.ReadFile(cp)
	if err != nil {
		return c, fmt.Errorf("readfile: %w", err)
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return c, fmt.Errorf("unmarshal: %w", err)
	}

	c.Theme = findTheme(root, c.Theme)
	klog.Infof("Config from %s: %+v", root, c)
	return c, nil
}
