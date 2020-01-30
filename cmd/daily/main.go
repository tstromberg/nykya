package main

import (
	"flag"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"k8s.io/klog"
)

var (
	rootFlag     string
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
		Name:  "daily",
		Usage: "daily mogger",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "root",
				Value:       ".",
				Usage:       "root directory to search for data",
				Destination: &rootFlag,
			},
		},

		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a post",
				Action:  func(c *cli.Context) error {
					return addCmd(c)
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "description",
						Value: "",
					},
				}
			},
			{
				Name:    "render",
				Aliases: []string{"r"},
				Usage:   "render posts",
				Action: func(c *cli.Context) error {
					return renderCmd(c)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
