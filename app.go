package main

import (
	"fmt"
	"github.com/kshitij10496/gftd/gftd"
	"github.com/urfave/cli"
	"os"
	"sort"
	"time"
)

func main() {
	initCommand := &cli.Command{
		Name:  "init",
		Usage: "Initializes the gftd application",
		Action: func(c *cli.Context) error {
			fmt.Printf("Initialized the application at: %s\n", gftd.UserHomeDir())
			return InitApp()
		},
	}

	newCommand := &cli.Command{
		Name:    "new",
		Aliases: []string{"add"},
		Usage:   "Add a new goal for today",
		Action: func(c *cli.Context) error {
			goal, err := PromptGoal()
			if err != nil {
				return fmt.Errorf("Unable to read your goal:", err)
			}

			werr := WriteGoal(goal)
			if werr != nil {
				return fmt.Errorf("Unable to write your goal:", werr)
			}

			fmt.Printf("Added the new goal: %+v\n", *goal)
			return nil
		},
		ArgsUsage: "[goal text]",
	}

	logCommand := &cli.Command{
		Name:  "log",
		Usage: "View your entire goal log",
		Before: func(c *cli.Context) error {
			exists, err := gftd.IsDBExists()
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

	achievedCommand := &cli.Command{
		Name:  "achieved",
		Usage: "Mark a goal as achieved",
		Action: func(c *cli.Context) error {
			fmt.Printf("You have successfully achieved your goal number: %v\n", c.Args().First())
			return nil
		},
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{*initCommand, *newCommand, *logCommand, *achievedCommand}

	app.Before = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "Welcome to GFTD!\n") // TODO: Put ASCII art
		return nil
	}

	app.After = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "Let's get working!\n") // TODO: Add a motivation quote
		return nil
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}

func ViewGoals() error {
	file, err := os.Open(gftd.DBFILE)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()

	goals, err := gftd.ReadAllGoals(file)
	if err != nil {
		return fmt.Errorf("Error fetching your goals: %v", err)
	}

	table := gftd.GetTableView(goals)
	fmt.Println(table)
	return nil
}

func InitApp() error {
	if exists, err := gftd.IsDBExists(); exists {
		if err == nil {
			fmt.Println("The application has already been initialized.")
			return nil
		}
	}

	if err := gftd.CreateDB(); err != nil {
		return fmt.Errorf("Error while setting up the database:", err)
	}

	return nil
}

func PromptGoal() (*gftd.Goal, error) {
	prompt := "What is your goal for today?"
	fmt.Println(prompt)
	message, err := gftd.ReadGoal(os.Stdin)
	if err != nil {
		return nil, err
	}

	return &gftd.Goal{message, time.Now(), false}, nil
}

func WriteGoal(goal *gftd.Goal) error {
	file, err := os.Open(gftd.DBFILE)
	if err != nil {
		return err
	}

	goals, err := gftd.ReadAllGoals(file)
	if err != nil {
		return err
	}
	file.Close()

	wfile, err := os.OpenFile(gftd.DBFILE, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer wfile.Close()

	goals = append(goals, *goal)
	return gftd.WriteAllGoals(wfile, goals)
}
