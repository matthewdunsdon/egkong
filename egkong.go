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

	lines := []string{"", "Examples:"}

	for i, ex := range examples.Examples {
		lines = append(lines, "  "+ex.Cli(examples.Context))
		lines = append(lines, "    "+ex.Description)
		if i < len(examples.Examples)-1 {
			lines = append(lines, "")
		}
	}

	for _, line := range lines {
		_, err := io.WriteString(ctx.Stdout, line+"\n")
		if err != nil {
			return err
		}
	}
	return err
}

func getCommand(ctx *kong.Context) (command string) {
	command = ""

	selected := ctx.Selected()
	if selected != nil {
		examples := selected.Tag.GetAll("examples")
		_ = examples
		command = strings.Replace(selected.FullPath(), ctx.Model.Name+" ", "", 1)
	}
	return
}

// New creates both kong and egcmd application, which have the help printer
// configured so that examples are shown.
func New(cli interface{}, options ...kong.Option) (parser *kong.Kong, appExamples *egcmd.App, err error) {
	options = append(options,
		kong.Help(func(o kong.HelpOptions, ctx *kong.Context) (err error) { return helpPrinter(appExamples, o, ctx) }),
	)
	parser, err = kong.New(cli, options...)
	appExamples = egcmd.New(parser.Model.Name)

	return
}
