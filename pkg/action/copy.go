package action

import (
	"fmt"
	"io"
	"os"

	"k8s.io/klog"
)

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
