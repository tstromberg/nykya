package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rjeczalik/notify"
	"github.com/tstromberg/nykya/pkg/nykya"
	"github.com/tstromberg/nykya/pkg/render"
	"github.com/tstromberg/nykya/pkg/store"
	"k8s.io/klog/v2"
)

type devCmd struct {
	Port   int  `default:32080 help:"Set a port TCP number"`
	Drafts bool `optional help:"Include draft posts"`
}

func renderLoop(ctx context.Context, dc nykya.Config) error {
	klog.Infof("starting render loop ...")
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(dc.In, c, notify.Remove); err != nil {
		klog.Fatal(err)
	}
	defer notify.Stop(c)

	for {
		ei := <-c
		klog.Infof("Got event:", ei)
		items, err := store.Scan(ctx, dc.Root)
		if err != nil {
			klog.Fatalf("scan: %w", err)
		}
		_, err = render.Site(ctx, dc, items)
		if err != nil {
			klog.Fatalf("render: %v", err)
		}
	}
}

func (c *devCmd) Run(globals *Globals) error {
	dc, err := nykya.ConfigFromRoot(globals.Root)
	if err != nil {
		return fmt.Errorf("config from root: %w", err)
	}

	if c.Drafts {
		dc.IncludeDrafts = true
	}

	ctx := context.Background()
	items, err := store.Scan(ctx, dc.In)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	paths, err := render.Site(ctx, dc, items)
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
