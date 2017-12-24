package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

func LogCommand() *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "View your entire goal log",
		Before: func(c *cli.Context) error {
			exists, err := IsDBExists()
			if !exists || err != nil {
				e := fmt.Errorf("You need to initialize the application using:\n $ gftd init\n")
				fmt.Println(e)
				return e // TODO: Find a way to disable help text
			}

			fmt.Fprintf(c.App.Writer, "Fetching your goals\n") // TODO: Add a progress bar
			return nil
		},
		Action: func(c *cli.Context) error {
			if err := ViewGoals(); err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	}
}

func ViewGoals() error {
	goals, err := GetGoals()
	if err != nil {
		return err
	}
	table := GetTableView(goals)
	fmt.Println(table)
	return nil
}
