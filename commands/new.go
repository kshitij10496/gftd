package commands

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
		Action: func(c *cli.Context) error {
			goal, err := PromptGoal()
			if err != nil {
				return fmt.Errorf("Unable to read your goal:", err)
			}

			if goal.Message == "" {
				fmt.Println("Don't be afraid to commit to your goal. You can do it.")
				return fmt.Errorf("Empty goal message")
			}

			werr := WriteGoal(goal)
			if werr != nil {
				return fmt.Errorf("Unable to write your goal:", werr)
			}

			fmt.Println("\nNow that you have committed to your goal, go get it!")
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

func WriteGoal(goal *Goal) error {
	goals, err := GetGoals()
	if err != nil {
		return nil
	}

	wfile, err := os.OpenFile(DBFILE, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer wfile.Close()

	goals = append(goals, goal)
	return WriteAllGoals(wfile, goals)
}

func GetGoals() ([]*Goal, error) {
	file, err := os.Open(DBFILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ReadAllGoals(file)
}
