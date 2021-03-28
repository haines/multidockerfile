package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
)

type context struct {
	Stdout io.Writer
	Stderr io.Writer
}

func Run() {
	new(
		os.Stdout,
		os.Stderr,
		os.Exit,
	).Run(os.Args[1:])
}

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func Capture(args ...string) (result Result) {
	exitCode := -1
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exited := errors.New("exited")

	defer func() {
		if exitCode != -1 {
			v := recover()
			if v != nil && v != exited {
				panic(fmt.Errorf("recovered from unexpected panic: %v", v))
			}
		}

		result = Result{
			ExitCode: exitCode,
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
		}
	}()

	new(
		&stdout,
		&stderr,
		func(code int) {
			exitCode = code
			panic(exited)
		},
	).Run(args)

	return
}

type command struct {
	parser *kong.Kong
}

func new(stdout io.Writer, stderr io.Writer, exit func(int)) *command {
	parser, err := kong.New(
		&rootCommand{},
		kong.Name("multidockerfile"),
		kong.Description("Split multi-stage Dockerfiles into multiple files."),
		kong.UsageOnError(),
		kong.Writers(stdout, stderr),
		kong.Exit(exit),
	)
	if err != nil {
		panic(err)
	}

	return &command{
		parser: parser,
	}
}

func (c *command) Run(args []string) {
	ctx, err := c.parser.Parse(args)
	c.parser.FatalIfErrorf(err)

	err = ctx.Run(&context{
		Stdout: ctx.Stdout,
		Stderr: ctx.Stderr,
	})
	ctx.FatalIfErrorf(err)

	ctx.Exit(0)
}
