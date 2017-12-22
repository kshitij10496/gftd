package main

import (
	"fmt"
	"github.com/kshitij10496/gftd/gftd"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"os"
	"sort"
	"strings"
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
			err := AchieveGoal()
			if err != nil {
				fmt.Println(err)
				return err
			}
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

func GetGoals() ([]*gftd.Goal, error) {
	file, err := os.Open(gftd.DBFILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return gftd.ReadAllGoals(file)
}

func ViewGoals() error {
	goals, err := GetGoals()
	if err != nil {
		return err
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
	goals, err := GetGoals()
	if err != nil {
		return nil
	}

	wfile, err := os.OpenFile(gftd.DBFILE, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer wfile.Close()

	goals = append(goals, goal)
	return gftd.WriteAllGoals(wfile, goals)
}

func AchieveGoal() error {
	goalPrompt := promptui.Prompt{Label: "Enter the goal"}
	goal, err := goalPrompt.Run()

	goal = strings.TrimSpace(strings.ToLower(goal))

	goalWords := strings.Fields(goal)
	goals, err := GetGoals()
	if err != nil {
		return err
	}

	strength := make(map[string]int)

	for _, g := range goals {
		for _, gw := range goalWords {
			if strings.Contains(strings.ToLower(g.Message), gw) {
				if _, found := strength[g.Message]; !found {
					strength[g.Message] = 0
				}
				strength[g.Message] += 1
			}
		}
	}

	fmt.Println("Total possible goals:", len(strength))
	if len(strength) == 0 {
		return fmt.Errorf("Could not find your goal. Enter a better message.")
	}

	pairs := rankByWordCount(strength)
	possibleGoals := []string{}

	for _, pair := range pairs {
		possibleGoals = append(possibleGoals, pair.Key)
	}

	selectGoal := promptui.Select{
		Label: "Select Goal",
		Items: possibleGoals,
		Size:  10,
	}

	i, s, e := selectGoal.Run()
	fmt.Println("I:", i)
	fmt.Println("S:", s)
	fmt.Println("E:", e)

	for _, g := range goals {
		if g.Message == s {
			g.Achieved = true
		}
	}

	file, err := os.OpenFile(gftd.DBFILE, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return gftd.WriteAllGoals(file, goals)
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
