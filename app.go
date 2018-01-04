package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/kshitij10496/gftd/cmd"
	"github.com/urfave/cli"
)

const LOGO = `
	        __ _      _
	  __ _ / _| |_ __| |
	 / _  | |_| __/ _  |
	| (_| |  _| || (_| |
	 \__/ |_|  \__\____|
	 |___/

`

func main() {
	initCommand := cmd.InitCommand()
	newCommand := cmd.NewCommand()
	logCommand := cmd.LogCommand()
	achieveCommand := cmd.AchieveCommand()

	app := cli.NewApp()
	app.Name = "gftd"
	app.Version = "0.1.0"
	app.HelpName = "gftd"
	app.Usage = "A tool to track your daily goals"
	app.Commands = []cli.Command{*initCommand, *newCommand, *logCommand, *achieveCommand}

	app.Before = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, LOGO)
		return nil
	}
	app.After = func(c *cli.Context) error {
		quote, err := cmd.GetMotivationalQuote()
		if err != nil {
			return err
		}
		fmt.Fprintf(c.App.Writer, quote)
		return nil
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
