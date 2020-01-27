package action

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"k8s.io/klog"
)

type AddOptions struct {
	Description string
	Timestamp   time.Time
}

func Add(path string, root string, opts AddOptions) error {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	klog.Infof("adding %s to %s, ts=%s", path, root, ts)
	y, m, d := ts.Date()
	dir := filepath.Join(root, fmt.Sprintf("%y-%m-%d", y, m, d))
	return copy(path, dir)
}

// copy copies a file.
func copy(src string, dst string) error {
	klog.Infof("copying %s -> %s", src, dst)
	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if _, err := io.Copy(d, s); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return d.Close()
}
