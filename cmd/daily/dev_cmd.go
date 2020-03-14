package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rjeczalik/notify"
	"github.com/tstromberg/daily/pkg/action"
	"github.com/tstromberg/daily/pkg/daily"
	"k8s.io/klog"
)

type DevCmd struct {
	Port int `help:"Set a port TCP number"`
}

func renderLoop(ctx context.Context, dc daily.Config) {
	klog.Infof("starting render loop ...")
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(dc.In, c, notify.Remove); err != nil {
		klog.Fatal(err)
	}
	defer notify.Stop(c)

	for {
		ei := <-c
		klog.Infof("Got event:", ei)
		_, err := action.Render(ctx, dc)
		if err != nil {
			klog.Fatalf("render: %v", err)
		}
	}
}

func (c *DevCmd) Run(globals *Globals) error {
	dc, err := daily.ConfigFromRoot(globals.Root)
	if err != nil {
		return fmt.Errorf("config from root: %w", err)
	}

	ctx := context.Background()
	paths, err := action.Render(ctx, dc)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	klog.Infof("rendered: %v", paths)

	klog.Infof("Starting up server on port %d ...", c.Port)
	fs := http.FileServer(http.Dir(dc.Out))
	http.Handle("/", fs)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil)
	return nil
}
