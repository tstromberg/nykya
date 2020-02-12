package action

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/tstromberg/daily/pkg/parse"
	"k8s.io/klog"
)

// AddOptions are options that can be passed to the add command
type AddOptions struct {
	Root        string
	Description string
	Timestamp   time.Time
}

// Add an object from a path
func Add(path string, opts AddOptions) error {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	klog.Infof("adding %s to %s, ts=%s opts=%+v", path, opts.Root, ts, opts)

	if path == "note" {
		return addNote(opts)
	}

	if _, err := url.Parse(path); err == nil {
		return addURL(path, opts)
	}

	if _, err := os.Stat(path); err == nil {
		return addFile(path, opts)
	}
	return nil
}

// addNote is for adding notes
func addNote(opts AddOptions) error {
	klog.Infof("addNote: %+v", opts)
	p := parse.Post{
		Description: opts.Description,
	}
	_ = guessDestination(opts.Root, p, opts)
	return nil
}

// addURL is for adding URL's
func addURL(path string, opts AddOptions) error {
	klog.Infof("addURL: %s - %+v", path, opts)
	_, err := url.Parse(path)
	if err != nil {
		return err
	}
	return nil
}

// addFile is for adding a local file
func addFile(path string, opts AddOptions) error {
	klog.Infof("addFile: %s - %+v", path, opts)
	return nil
}

// defaultHierarchy returns an unconfigured hierarchy. Deal with it.
func defaultHierarchy(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%d-%d", y, m, d)
}

// guessDestination picks the local destination of the file
func guessDestination(root string, p parse.Post, opts AddOptions) string {
	return filepath.Join(root, defaultHierarchy(p.Posted), filepath.Base(p.Source))
}
