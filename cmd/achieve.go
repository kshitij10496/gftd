package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

func AchieveCommand() *cli.Command {
	return &cli.Command{
		Name:  "achieve",
		Usage: "Mark a goal as achieved",
		Before: func(c *cli.Context) error {
			exists, err := IsDBExists()
			if !exists || err != nil {
				e := fmt.Errorf("You need to initialize the application using:\n $ gftd init\n")
				fmt.Println(e)
				return e // TODO: Find a way to disable help text
			}
			// TODO: Check for number of goals in the database
			return nil
		},
		Action: func(c *cli.Context) error {
			err := AchieveGoal()
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Printf("You have achieved an important goal. More power to you!\n")
			return nil
		},
	}
}

func AchieveGoal() error {
	// TODO: Check if any goals exist before prompting
	goalPrompt := promptui.Prompt{Label: "Enter the goal"}
	goal, err := goalPrompt.Run()

	goal = strings.TrimSpace(strings.ToLower(goal))

	goalWords := strings.Fields(goal)
	goals, err := ReadAllGoals()
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

	_, s, e := selectGoal.Run()
	if e != nil {
		return e
	}

	return UpdateGoal(s)
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
