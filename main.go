package main

import (
	//"bufio"
	//"encoding/json"
	"fmt"
	"os"
	//"strings"
	"flag"
	"github.com/kshitij10496/gftd/gftd"
	"time"
)

// func main() {

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
	flag.Parse()

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
		fmt.Println("Number of goals:", len(goals))
		for i, goal := range goals {
			fmt.Printf("%d. %v\n", i+1, goal)
		}
		return
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
}
