package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initializes the gftd application",
		Action: func(c *cli.Context) error {
			fmt.Printf("Initialized the application at: %s\n", UserHomeDir())
			return InitAction()
		},
	}
}

func InitAction() error {
	if exists, err := IsDBExists(); exists {
		if err == nil {
			fmt.Println("The application has already been initialized.")
			return nil
		}
	}

	if err := CreateDB(); err != nil {
		return fmt.Errorf("Error while setting up the database:", err)
	}

	return nil
}
