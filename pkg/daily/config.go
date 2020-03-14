package daily

import (
	"path/filepath"

	"k8s.io/klog"
)

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
	c := Config{
		Root: root,
		In:   filepath.Join(root, "in"),
		Out:  filepath.Join(root, "out"),

		Organization: map[string]string{
			"thought": "thoughts/{{.Year}}-{{.Month}}-{{.Day}}/{{ .Index }}-{{ .Slug }}",
		},
	}
	klog.Infof("Config from %s: %+v", root, c)
	return c, nil
}
