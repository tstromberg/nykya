package main

import (
	"flag"

	"github.com/alecthomas/kong"
	"k8s.io/klog"
)

// Globals are global flags that can be set for all subcommands
type Globals struct {
	Root string `help:"Set the root directory for the site"`
}

// CLI defines the subcommands and flags supported
type CLI struct {
	Globals
	Add    addCmd    `cmd:"" help:"Add content"`
	Render renderCmd `cmd:"" help:"Render output to a static directory"`
	Dev    devCmd    `cmd:"" help:"Developer mode: local webserver with instant rendering"`
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("nykya"),
		kong.Description("nykya personal site generator"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
