package store

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/jinzhu/copier"
	"github.com/tstromberg/nykya/pkg/nykya"
)

func editorCmd(ctx context.Context, dc nykya.Config, goos string, path string) *exec.Cmd {
	if goos == "windows" {
		return exec.Command("cmd", "/C", "start", path)
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	return exec.Command(editor, path)
}

func openEditor(ctx context.Context, dc nykya.Config, i nykya.RawItem) (nykya.RawItem, error) {
	ni := nykya.RawItem{}
	copier.Copy(i, ni)

	tf, err := ioutil.TempFile("", fmt.Sprintf("*%s", extForFormat(i.Format)))
	if err != nil {
		return ni, fmt.Errorf("tempfile: %w", err)
	}
	if err := saveRawItem(ctx, dc, i, tf.Name()); err != nil {
		return ni, fmt.Errorf("saveRawItem: %w", err)
	}
	tf.Close()

	c := editorCmd(ctx, dc, runtime.GOOS, tf.Name())
	out, err := c.CombinedOutput()
	if err != nil {
		return ni, fmt.Errorf("run %v: %w\n%s", c.Args, err, out)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\nPress enter to save, Ctrl-C to abort -> ")
	if _, err := reader.ReadString('\n'); err != nil {
		return ni, fmt.Errorf("readstring: %w", err)
	}

	fm, err := fromMarkdown(tf.Name())
	return *fm, err
}
