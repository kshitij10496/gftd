package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initializes the gftd application",
		Action: func(c *cli.Context) error {
			err := InitAction()
			if err != nil {
				return err
			}

			fmt.Printf("Initialized the application at: %s\n", DBFILE)
			return nil
		},
	}
}

func InitAction() error {
	if exists, err := IsDBExists(); exists {
		if err == nil {
			e := fmt.Errorf("The application has already been initialized.")
			fmt.Println(e)
			return e
		}
	}

	if err := CreateDB(); err != nil {
		e := fmt.Errorf("Error while setting up the database:", err)
		fmt.Println(e)
		return e
	}
	return nil
}
