package nykya

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v1"
	"k8s.io/klog/v2"
)

// Configuration file name
const configFileName = "nykya.yaml"

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
	}
	return ""
}

func defaultConfig(root string) Config {
	return Config{
		Root:        root,
		In:          root,
		Out:         root + ".out",
		Title:       "Example Title",
		Subtitle:    "Example Subtitle",
		Description: "Example Description",
		Theme:       "basic",
	}
}

// ConfigFromRoot returns the sites configuration
func ConfigFromRoot(rootOverride string) (Config, error) {
	root := "."
	envRoot := os.Getenv("NYKYA_ROOT")

	if rootOverride != "" {
		klog.Infof("root set to: %v", rootOverride)
		root = rootOverride
	} else if envRoot != "" {
		klog.Infof("NYKYA_ROOT set to: %v", envRoot)
		root = envRoot
	}

	cp := filepath.Join(root, configFileName)
	if _, err := os.Stat(cp); err != nil {
		return Config{}, fmt.Errorf("Unable to find %s within --root (%q), $NYKYA_ROOT (%q), and the current directory", configFileName, rootOverride, envRoot)
	}

	b, err := ioutil.ReadFile(cp)
	if err != nil {
		return Config{}, fmt.Errorf("readfile: %w", err)
	}

	root, err = filepath.Abs(root)
	if err != nil {
		return Config{}, fmt.Errorf("abs: %w", err)
	}

	c := defaultConfig(root)
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return c, fmt.Errorf("unmarshal: %w", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return c, fmt.Errorf("getwd: %w", err)
	}

	if err := os.Chdir(root); err != nil {
		return c, fmt.Errorf("chdir: %w", err)
	}

	in, err := filepath.Abs(filepath.FromSlash(c.In))
	if err != nil {
		return c, fmt.Errorf("abs: %w", err)
	}
	c.In = in

	out, err := filepath.Abs(filepath.FromSlash(c.Out))
	if err != nil {
		return c, fmt.Errorf("abs: %w", err)
	}
	c.Out = out

	if err := os.Chdir(root); err != nil {
		return c, fmt.Errorf("chdir: %w", err)
	}

	if err := os.Chdir(pwd); err != nil {
		return c, fmt.Errorf("chdir: %w", err)
	}

	c.Theme = findTheme(root, c.Theme)
	klog.Infof("Config from %s: %+v", root, c)
	return c, nil
}
