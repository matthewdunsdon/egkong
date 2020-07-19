package egkong

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/kong"

	"github.com/matthewdunsdon/egcmd"
)

type testExamplesFinder struct {
	lastCommand string
}

func (t *testExamplesFinder) Find(command string) (examples egcmd.ExamplesFound) {
	t.lastCommand = command
	return
}

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

func TestNew(t *testing.T) {
	os.Args[0] = "os-app-name"
	var (
		parser, appExamples = New(&cli)
		expectedName        = "os-app-name"
	)

	if got := parser.Model.Name; expectedName != got {
		t.Errorf("Expected app name to be %q, got %q", expectedName, got)
	}

	if got := appExamples.Find("").Context; expectedName != got {
		t.Errorf("Expected app examples name to be %q, got %q", expectedName, got)
	}
}

func TestNewWithNameOptionSpecified(t *testing.T) {
	var (
		parser, appExamples = New(&cli, kong.Name("app-name"))
		expectedName        = "app-name"
	)

	if got := parser.Model.Name; expectedName != got {
		t.Errorf("Expected app name to be %q, got %q", expectedName, got)
	}

	if got := appExamples.Find("").Context; expectedName != got {
		t.Errorf("Expected app examples name to be %q, got %q", expectedName, got)
	}
}

func TestGetCommand(t *testing.T) {

	testCases := []struct {
		testName string
		args     []string
		want     string
	}{
		{
			"ForApp",
			[]string{"--help"},
			"",
		},
		{
			"ForCommand",
			[]string{"init", "--help"},
			"init",
		},
		{
			"ForSubCommand",
			[]string{"config", "snapshot", "--help"},
			"config snapshot",
		},
		{
			"ForPartOfSubCommand",
			[]string{"config", "--help"},
			"config",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			helpWasCalled, _, _ := parseWithHelpPrinter(tc.args, func(o kong.HelpOptions, ctx *kong.Context) (err error) {
				actual := getCommand(ctx)
				if got := actual; tc.want != got {
					t.Errorf("Expected getCommand to return %#v, got %#v", tc.want, actual)
				}
				return nil
			})

			if !helpWasCalled {
				t.Errorf("Expected help function to be called")
			}
		})
	}
}

func TestHelpPrinter(t *testing.T) {
	examples := egcmd.New("app-name")
	examples.Example("init", "Ius legimus nonumes te, pri dicat nominavi copiosae id, odio rebum facilis ea pro.")
	examples.Example("config snapshot odio", "At vis primis debitis, ei verear omittantur.")
	examples.Command("version").Example("--json", "At vis primis debitis, ei verear omittantur.")
	examples.Command("config").Example("snapshot odio", "At vis primis debitis, ei verear omittantur.")
	examples.Command("config snapshot").Example("odio", "At vis primis debitis, ei verear omittantur.")

	testCases := []struct {
		testName     string
		args         []string
		wantExamples string
	}{
		{
			"ForApp",
			[]string{"--help"},
			"\n\nExamples:\n  app-name init\n    Ius legimus nonumes te, pri dicat nominavi copiosae id, odio rebum facilis ea pro.\n\n  app-name config snapshot odio\n    At vis primis debitis, ei verear omittantur.\n",
		},
		{
			"ForCommand",
			[]string{"version", "--help"},
			"\n\nExamples:\n  app-name version --json\n    At vis primis debitis, ei verear omittantur.\n",
		},
		{
			"ForCommandWithoutExamples",
			[]string{"init", "--help"},
			"",
		},
		{
			"ForSubCommand",
			[]string{"config", "snapshot", "--help"},
			"\n\nExamples:\n  app-name config snapshot odio\n    At vis primis debitis, ei verear omittantur.\n",
		},
		{
			"ForPartOfSubCommand",
			[]string{"config", "--help"},
			"\n\nExamples:\n  app-name config snapshot odio\n    At vis primis debitis, ei verear omittantur.\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			helpWasCalled, stdout, stderr := parseWithHelpPrinter(tc.args, func(o kong.HelpOptions, ctx *kong.Context) (err error) {
				err = helpPrinter(examples, o, ctx)
				return
			})

			if !helpWasCalled {
				t.Errorf("Expected help function to be called")
			}

			examplesFound := strings.Contains(stdout, "\n\nExamples:\n")
			expectedExamples := tc.wantExamples != ""

			if expectedExamples {
				if !examplesFound {
					t.Errorf("Expected stdout to have examples, got %q", stdout)
				}

				if got := stdout; !strings.HasSuffix(got, tc.wantExamples) {
					t.Errorf("Expected stdout output end with %q, got %q", tc.wantExamples, got)
				}
			} else if !expectedExamples && examplesFound {
				t.Errorf("Expected stdout not to have examples, got %q", stdout)
			}

			if got := stderr; "" != got {
				t.Errorf("Expected stderr output to be empty, got %q", got)
			}
		})
	}
}

func parseWithHelpPrinter(args []string, help kong.HelpPrinter) (helpWasCalled bool, stdout, stderr string) {
	var stdoutWriter, stderrWriter bytes.Buffer

	helpWasCalled = false
	parser, _ := kong.New(&cli,
		kong.Name("app-name"),
		kong.Exit(func(res int) {}),
		kong.Writers(&stdoutWriter, &stderrWriter),
		kong.Help(func(o kong.HelpOptions, ctx *kong.Context) error {
			helpWasCalled = true
			return help(o, ctx)
		}),
	)

	parser.Parse(args)
	stdout = stdoutWriter.String()
	stderr = stderrWriter.String()

	return
}
