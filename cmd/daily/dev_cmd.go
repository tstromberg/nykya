package main

import (
	"fmt"
	"net/http"

	"github.com/rjeczalik/notify"
	"github.com/tstromberg/daily/pkg/action"
	"github.com/tstromberg/daily/pkg/daily"
	"k8s.io/klog"
)

func renderLoop(dc daily.Config) {
	klog.Infof("starting render loop ...")
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(dc.In, c, notify.Remove); err != nil {
		klog.Fatal(err)
	}
	defer notify.Stop(c)

	for {
		ei := <-c
		klog.Infof("Got event:", ei)
		_, err := action.Render(dc)
		if err != nil {
			klog.Fatalf("render: %v", err)
		}
	}
}

func devCmd(root string, port int) {
	dc := daily.ConfigFromRoot(root)
	_, err := action.Render(dc)
	if err != nil {
		klog.Fatalf("render: %v", err)
	}

	klog.Infof("Starting up server on port %d ...", port)
	fs := http.FileServer(http.Dir(dc.Out))
	http.Handle("/", fs)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
