package main

import (
	"github.com/Forget-C/gopkg/cmd/executor/app"
	"github.com/Forget-C/gopkg/pkg/content/cli"
	"github.com/common-nighthawk/go-figure"
	"os"
)

func main() {
	myFigure := figure.NewFigure("executor", "", true)
	myFigure.Print()
	command := app.NewPreInitCommand()
	code := cli.Run(command)
	os.Exit(code)
}
