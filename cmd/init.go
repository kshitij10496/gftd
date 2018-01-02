package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initializes the gftd application",
		Before: func(c *cli.Context) error {
			if exists, err := IsDBExists(); exists || err != nil {
				e := fmt.Errorf("The application has already been initialized.")
				fmt.Println(e)
				return e // TODO: Do not show the help text
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			err := InitAction()
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Printf("Initialized the application at: %s\n", DBFILE)
			return nil
		},
	}
}

func InitAction() error {
	if err := CreateDB(); err != nil {
		e := fmt.Errorf("Error while setting up the database: %v", err)
		return e
	}
	return nil
}
