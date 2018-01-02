package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:    "new",
		Aliases: []string{"add"},
		Usage:   "Add a new goal for today",
		Before: func(c *cli.Context) error {
			exists, err := IsDBExists()
			if !exists || err != nil {
				e := fmt.Errorf("You need to initialize the application using:\n $ gftd init\n")
				RED.Println(e)
				return e // TODO: Find a way to disable help text
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			goal, err := PromptGoal()
			if err != nil {
				return fmt.Errorf("Unable to read your goal:", err)
			}

			if goal.Message == "" {
				RED.Println("Don't be afraid to commit to your goal. You can do it.")
				return fmt.Errorf("Empty goal message")
			}

			werr := WriteGoal(goal)
			if werr != nil {
				er := fmt.Errorf("Unable to write your goal:", werr)
				RED.Println(er)
				return er
			}

			GREEN.Println("\nNow that you have committed to your goal, go get it!")
			return nil
		},
		ArgsUsage: "[goal text]",
	}
}

func PromptGoal() (*Goal, error) {
	prompt := fmt.Sprintf("%s: %s", "gftd", "What is your goal for today?")
	fmt.Println(prompt)
	fmt.Printf("%4s: ", "me")
	message, err := ReadGoal(os.Stdin)
	if err != nil {
		return nil, err
	}

	return &Goal{message, time.Now(), false}, nil
}
