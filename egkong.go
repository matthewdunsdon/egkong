package egkong

import (
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/matthewdunsdon/egcmd"
)

// ExamplesFinder is the interface that wraps the Find method.
//
// Find returns the examples that belong to a particular command.
type ExamplesFinder interface {
	Find(command string) egcmd.ExamplesFound
}

func helpPrinter(finder ExamplesFinder, options kong.HelpOptions, ctx *kong.Context) (err error) {
	examples := finder.Find(getCommand(ctx))

	err = kong.DefaultHelpPrinter(options, ctx)
	if err != nil || examples.Examples == nil {
		return
	}

	_, err = io.WriteString(ctx.Stdout, "\nExamples:\n")
	for i, ex := range examples.Examples {
		_, err = io.WriteString(ctx.Stdout, "  "+ex.Cli(examples.Context)+"\n")
		_, err = io.WriteString(ctx.Stdout, "    "+ex.Description+"\n")
		if i < len(examples.Examples)-1 {
			_, err = io.WriteString(ctx.Stdout, "\n")
		}
	}
	return err
}

func getCommand(ctx *kong.Context) (command string) {
	command = ""

	selected := ctx.Selected()
	if selected != nil {
		command = strings.Replace(selected.FullPath(), ctx.Model.Name+" ", "", 1)
	}
	return
}

// New creates both kong and egcmd application, which have the help printer
// configured so that examples are shown.
func New(cli interface{}, options ...kong.Option) (parser *kong.Kong, appExamples *egcmd.App) {
	options = append(options,
		kong.Help(func(o kong.HelpOptions, ctx *kong.Context) (err error) { return helpPrinter(appExamples, o, ctx) }),
	)
	parser, err := kong.New(cli, options...)
	appExamples = egcmd.New(parser.Model.Name)

	if err != nil {
		panic(err)
	}
	return
}
