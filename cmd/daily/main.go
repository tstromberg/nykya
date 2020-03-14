package main

import (
	"flag"

	"github.com/alecthomas/kong"
	"k8s.io/klog"
)

type Globals struct {
	Root string `help:"Set the debug directory"`
}

type CLI struct {
	Globals

	Add    AddCmd    `cmd:"" help:"Add files."`
	Render RenderCmd `cmd:"" help:"Render output."`
	Dev    DevCmd    `cmd:"" help:"Developer mode"`
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("daily"),
		kong.Description("daily mogger."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
