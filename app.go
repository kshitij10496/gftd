package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/kshitij10496/gftd/commands"
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
	initCommand := commands.InitCommand()
	newCommand := commands.NewCommand()
	logCommand := commands.LogCommand()
	achieveCommand := commands.AchieveCommand()

	app := cli.NewApp()
	app.Name = "gftd"
	app.Version = "0.0.1"
	app.HelpName = "gftd"
	app.Usage = "A tool to track your daily goals"
	//app.Description = "Your daily goal planner"
	app.Commands = []cli.Command{*initCommand, *newCommand, *logCommand, *achieveCommand}

	app.Before = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, LOGO)
		return nil
	}

	app.After = func(c *cli.Context) error {
		// fmt.Fprintf(c.App.Writer, "Let's get working!\n") TODO: Add a motivation quote
		return nil
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
