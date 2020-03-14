package parse

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/tstromberg/daily/pkg/daily"
	"gopkg.in/yaml.v1"
	"k8s.io/klog"
)

func fromYAML(path string) (*daily.Item, error) {
	klog.Infof("yaml: %s", path)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var i daily.Item
	err = yaml.Unmarshal(b, &i)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	// TODO: Find a more elegant way to handle front-matter
	si := bytes.Index(b, []byte(daily.DocumentSeparator))
	if si > 0 {
		i.SetContent(string(b[si+len(daily.DocumentSeparator):]))
	}

	klog.Infof("read: %+v", i)
	return &i, nil
}
