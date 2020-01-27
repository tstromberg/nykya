package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/urfave/cli/v2"
	"k8s.io/klog"
)

var (
	thumbQuality = 85
	indexTmpl    = template.Must(template.ParseFiles("index.tmpl"))
)

type Post struct {
	Image string
}

type Page struct {
	Title string
	Posts []*Post
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	app := &cli.App{
		Name:   "daily",
		Usage:  "daily mogger",
		Action: doMain,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func doMain(c *cli.Context) error {
	path := c.Args().Get(0)
	klog.Infof("image: %q", path)

	img, err := imgio.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	thumb := transform.Resize(img, 800, 800, transform.Linear)
	klog.Infof("writing to output.jpg")

	if err := imgio.Save("output.jpg", thumb, imgio.JPEGEncoder(thumbQuality)); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	p := Page{
		Title: "my title",
		Posts: []*Post{
			&Post{
				Image: "output.jpg",
			},
		},
	}
	indexTmpl.Execute(os.Stdout, p)
	return nil
}
