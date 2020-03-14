package action

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tstromberg/daily/pkg/daily"
	"gopkg.in/yaml.v1"
	"k8s.io/klog"
)

// AddOptions are options that can be passed to the add command
type AddOptions struct {
	Root string

	Title string
	Text  string
	Kind  string

	Timestamp time.Time
}

// Add an object from a path
func Add(ctx context.Context, dc daily.Config, opts AddOptions) error {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	klog.Infof("adding %q to %s, ts=%s opts=%+v", opts.Text, opts.Root, ts, opts)

	switch opts.Kind {
	case "thought":
		return addThought(ctx, dc, opts)
	default:
		return fmt.Errorf("unknown object type: %q", opts.Kind)
	}
}

// addThought is for adding thoughts
func addThought(ctx context.Context, dc daily.Config, opts AddOptions) error {
	klog.Infof("addNote %+v", opts)

	words := strings.Split(opts.Text, " ")
	slug := strings.Join(words[0:3], "-")

	i := daily.Item{
		Kind:    opts.Kind,
		Text:    opts.Text,
		Posted:  &opts.Timestamp,
		Created: &opts.Timestamp,
		Updated: &opts.Timestamp,
		Slug:    slug,
	}

	i.Hier = calculateHierarchy(dc, i)
	return saveItem(ctx, dc, i)
}

// saveItem saves an item to disk
func saveItem(ctx context.Context, dc daily.Config, i daily.Item) error {
	klog.Infof("marshalling: %+v", i)
	b, err := yaml.Marshal(i)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	path := itemPath(dc, i)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		klog.Infof("Creating %s ...", dir)
		err := os.MkdirAll(dir, 0600)
		if err != nil {
			klog.Errorf("mkdir(%s) failed: %v", dir, err)
			// Keep on truckin!
		}
	}

	fmt.Printf("Writing to %s ...", path)
	return ioutil.WriteFile(path, b, 0600)
}

func itemPath(dc daily.Config, i daily.Item) string {
	// TODO: respect input directory
	return filepath.Join(dc.Root, "in", i.Hier, i.Slug+".yaml")
}

// defaultHierarchy returns an unconfigured hierarchy. Deal with it.
func calculateHierarchy(dc daily.Config, i daily.Item) string {
	y, m, d := time.Now().Date()

	// TODO: respect site organization
	// NOTE: Unix paths here!
	return path.Join(i.Kind+"s", fmt.Sprintf("%d-%d-%d", y, m, d))
}
