package main

import (
	"flag"
	"fmt"

	"github.com/alecthomas/kong"
	"k8s.io/klog"
)

type addOpts struct {
	Description string   `help:"Set a description for the post"`
	Paths       []string `arg:"" optional:"" help:"Paths to add." type:"path"`
}

type renderOpts struct {
}

var cli struct {
	Root   string     `help:"Set the debug directory"`
	Add    addOpts    `cmd:"" help:"Add files."`
	Render renderOpts `cmd:"" help:"Render output."`
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	ctx := kong.Parse(&cli,
		kong.Name("daily"),
		kong.Description("daily mogger."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	switch ctx.Command() {
	case "add":
		addWithoutPath(cli.Root, cli.Add)
	case "add <paths>":
		addPaths(cli.Root, cli.Add)
	case "render":
		renderCmd(cli.Root)
	default:
		fmt.Printf("unknown command: %q\n", ctx.Command())
	}
}
