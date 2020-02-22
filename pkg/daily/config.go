package daily

import "path/filepath"

type Config struct {
	Title       string
	SubTitle    string
	Description string

	In  string
	Out string
}

func ConfigFromRoot(root string) Config {
	// TODO: Parse YAML file from root
	return Config{
		In:  filepath.Join(root, "in"),
		Out: filepath.Join(root, "out"),
	}
}
