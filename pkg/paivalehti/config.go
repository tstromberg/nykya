package paivalehti

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

	In  string
	Out string

	Organization map[string]string
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

	cp := filepath.Join(root, "paivalehti.yaml")
	if _, err := os.Stat(cp); err != nil {
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
	klog.Infof("Config from %s: %+v", root, c)
	return c, nil
}
