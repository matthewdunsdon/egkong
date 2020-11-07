package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/matthewdunsdon/egcmd/egkong"
)

var cli struct {
	Init struct {
	} `cmd:"init" help:"Initialise app data."`

	Version struct {
		JSON bool `help:"Get version details in json format."`
	} `cmd:"version" help:"Get application version details."`

	Config struct {
		Snapshot struct {
			Name string `arg`
		} `cmd`
	} `cmd`
}

func version(ctx *kong.Context) {
	fmt.Println("", cli.Version)
}

var (
	parser, appExamples = egkong.New(&cli, kong.Name("myapp"), kong.Description("This is my app."))
	_                   = appExamples.Example("init", "Ius legimus nonumes te, pri dicat nominavi copiosae id, odio rebum facilis ea pro.")

	versionCmdEx = appExamples.Command("version")
	_            = versionCmdEx.Example("--json", "At vis primis debitis, ei verear omittantur signiferumque mei, quo esse aperiri an. Dolore vocent consequuntur pro an, nam no iusto tamquam suscipit.")
	_            = versionCmdEx.Example("--json2", "At vis primis debitis, ei verear omittantur signiferumque mei, quo esse aperiri an. Dolore vocent consequuntur pro an, nam no iusto tamquam suscipit.")
	_            = versionCmdEx.Example("--json3", "At vis primis debitis, ei verear omittantur signiferumque mei, quo esse aperiri an. Dolore vocent consequuntur pro an, nam no iusto tamquam suscipit.")
)

func main() {
	app, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)
	switch app.Command() {
	case "rm <path>":
	case "version":
		version(app)
	default:
		panic(app.Command())
	}
}
