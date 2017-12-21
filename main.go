package main

import (
	"flag"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/kshitij10496/gftd/gftd"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
	"time"
)

// Hacky tabular representation
//
func GetTableView(goals []gftd.Goal) *uitable.Table {
	table := uitable.New()
	table.MaxColWidth = 50
	table.Wrap = true
	table.Separator = " | "

	// TODO: Find better ways to format
	table.AddRow("S.No", "Date", "Goal", "Achieved")
	table.AddRow("====", strings.Repeat("=", 16), strings.Repeat("=", 50), "========")
	for i, goal := range goals {
		year, month, day := goal.Timestamp.Date()
		table.AddRow(i+1, fmt.Sprintf("%d %v %d", day, month, year), goal.Message, goal.Achieved)
	}
	return table
}

func main() {
	// 1. Check if database exists
	//      a) No: Create Database and proceed
	// 2. Check the database for unchecked entries
	//      a) No: Prompt for new entry
	//      b) Exists: Fetch the latest unchecked entry in the database
	//                 Compare dates of the unchecked entry and current date
	//                 i> Same Date: Print the entry and exit
	//                 ii> Different Date: Prompt to check it off

	view := flag.Bool("view", false, "Lists all the goals")
	//achieved := flag.Bool("a", false, "Mark goals you have achieved")
	flag.Parse()

	fmt.Println(strings.Repeat("=", 79))
	file, err := os.Open(gftd.DBFILE)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()
	goals, err := gftd.ReadAllGoals(file)

	if *view {
		fmt.Println("Reading all the goals")
		if err != nil {
			fmt.Println("Unable to read:", err)
			os.Exit(1)
		}
		table := GetTableView(goals)
		fmt.Println(table)
		return
	}

	fmt.Println()

	previousGoal := &goals[len(goals)-1]
	if !previousGoal.Achieved {
		year, month, day := previousGoal.Timestamp.Date()
		previousGoalDate := fmt.Sprintf("%d %v %d", day, month, year)

		l := fmt.Sprintf("Your goal for %s was:", previousGoalDate)

		label := strings.Join([]string{l, previousGoal.Message}, "\n")
		fmt.Println(label)
		fmt.Println()

		prompt := promptui.Prompt{
			Label:     "Have you achieved it",
			IsConfirm: true,
		}

		//fmt.Println(previousGoal)
		_, err := prompt.Run()

		if err != nil {
			fmt.Printf("Have confidence in yourself. You can do it!\n")

		} else {
			//fmt.Println(result)
			fmt.Println("Kudos! More power to you for achieving your goal.")
			previousGoal.Achieved = true
		}
	}

	prompt := "What is your goal for today?"
	fmt.Println(prompt)
	message, err := gftd.ReadGoal(os.Stdin)
	if err != nil {
		fmt.Println("Unable to read your goal:", err)
		os.Exit(1)
	}

	goal := gftd.Goal{message, time.Now(), false}
	goals = append(goals, goal)

	// Hack for appending to the database
	fileWrite, err := os.OpenFile(gftd.DBFILE, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
	defer fileWrite.Close()

	werr := gftd.WriteAllGoals(fileWrite, goals)
	if werr != nil {
		fmt.Println("Unable to write your goal to the database:", werr)

		os.Exit(1)
	}
	fmt.Println("Successfully saved your goal")
	fmt.Println(strings.Repeat("=", 80))
}
