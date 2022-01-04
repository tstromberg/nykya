package nykya

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

// Configuration file name
const ConfigFileName = "nykya.yaml"

// DefaultOrganization shows where to put files if organization is unset
const DefaultOrganization = `{{ .Kind }}s/{{ .Date.Format "2006" }}`

// Config is site configuration
type Config struct {
	Root string

	Title       string
	Subtitle    string
	Description string

	SyncCommand string

	In    string
	Out   string
	Theme string

	IncludeDrafts bool

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
		In:          "content/",
		Out:         "rendered/",
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

	c := defaultConfig(root)
	cp := filepath.Join(root, ConfigFileName)
	if _, err := os.Stat(cp); err == nil {

		b, err := ioutil.ReadFile(cp)
		if err != nil {
			return Config{}, fmt.Errorf("readfile: %w", err)
		}

		root, err = filepath.Abs(root)
		if err != nil {
			return Config{}, fmt.Errorf("abs: %w", err)
		}

		err = yaml.Unmarshal(b, &c)
		if err != nil {
			return c, fmt.Errorf("unmarshal: %w", err)
		}
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
