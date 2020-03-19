package daily

import (
	"fmt"
	"path/filepath"

	"k8s.io/klog"
)

// If organization is unset, use this!
const DefaultOrganization = `{{ .Kind }}s/{{ .Posted.Format "2006-01-02" }}`

// Config is site configuration
type Config struct {
	Root string

	Title       string
	SubTitle    string
	Description string

	In  string
	Out string

	Organization map[string]string
}

// ConfigFromRoot returns the sites configuration
func ConfigFromRoot(root string) (Config, error) {
	// TODO: Parse YAML file from root
	root, err := filepath.Abs(root)
	if err != nil {
		return Config{}, fmt.Errorf("abs: %w", err)
	}

	c := Config{
		Root: root,
		In:   filepath.Join(root, "in"),
		Out:  filepath.Join(root, "out"),
	}
	klog.Infof("Config from %s: %+v", root, c)
	return c, nil
}
