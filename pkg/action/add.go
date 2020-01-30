package action

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"k8s.io/klog"
)

type AddOptions struct {
	Description string
	Timestamp   time.Time
}

// defaultHierarchy returns an unconfigured hierarchy. Deal with it.
func defaultHierarchy() {
	y, m, d := ts.Date()
	return fmt.Sprintf("%y-%m-%d", y, m, d)
}

func Add(path string, root string, opts AddOptions) error {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	klog.Infof("adding %s to %s, ts=%s opts=%+v", path, root, ts, opts)

	if path == "note" {
		return addNote(root, opts)
	}

	if _, err := url.Parse(s); err == nil {
		return addURL(path, root, opts)
	}

	if _, err := os.Stat(path); err == nil {
		return addFile(path, root, opts)
	}
}

func addNote()

func addURL() {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
}
