package gftd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gosuri/uitable"
)

type Goal struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Achieved  bool      `json:"achieved"`
}

var DBFILE = UserHomeDir() + "/.gftd.json"

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// Create a database file
//
func CreateDB() error {
	file, err := os.Create(DBFILE)
	defer file.Close()

	return err
}

// Checks if the database exists or not.
//
func IsDBExists() (bool, error) {
	if _, err := os.Stat(DBFILE); err != nil {
		if os.IsNotExist(err) {
			return false, nil // file does not exist
		} else if os.IsPermission(err) {
			return true, err // file exists but permission denied
		} else {
			return true, err // file exits but some other error
			panic(err)
		}
	}
	return true, nil
}

// Reads all the goals currently in the database
//
func ReadAllGoals(file io.Reader) ([]*Goal, error) {
	var goals []*Goal
	dec := json.NewDecoder(file)
	if err := dec.Decode(&goals); err != nil && err != io.EOF {
		return nil, err // TODO: Handle the case of empty file EOF error
	}

	return goals, nil
}

// Writes a goal to the database
//
func WriteAllGoals(file io.Writer, goals []*Goal) error {
	enc := json.NewEncoder(file)
	return enc.Encode(&goals)
}

// Reads a single-line goal from a reader.
//
func ReadGoal(r io.Reader) (string, error) {
	buffReader := bufio.NewReader(r)
	goal, err := buffReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	goal = strings.Replace(goal, "\n", "", -1)
	return goal, nil
}

// Encoded a goal.
//
func encodeGoal(g *Goal) ([]byte, error) {
	b, err := json.Marshal(*g)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Decode a goal.
//
func decodeGoal(encodedGoal []byte) (*Goal, error) {
	g := new(Goal)
	err := json.Unmarshal(encodedGoal, g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Checks if both the time variables have the same date.
//
func isSameDate(time1, time2 time.Time) bool {
	year1, month1, day1 := time1.Date()
	year2, month2, day2 := time2.Date()

	if year1 == year2 && month1 == month2 && day1 == day2 {
		fmt.Println("It's the same date")
		return true
	}
	return false
}

// Hacky tabular representation
//
func GetTableView(goals []*Goal) *uitable.Table {
	table := uitable.New()
	table.MaxColWidth = 50
	table.Wrap = true
	table.Separator = " | "

	// TODO: Find better ways to format
	table.AddRow("S.No", "Date", "Goal", "Achieved")
	table.AddRow("====", strings.Repeat("=", 16), strings.Repeat("=", 50), "========")
	for i, goal := range goals {
		year, month, day := goal.Timestamp.Date()
		goalStatus := func(achieved bool) string {
			status := "NO"
			if achieved {
				status = "YES"
			}
			return fmt.Sprintf("%3s", status)
		}
		table.AddRow(i+1, fmt.Sprintf("%d %v %d", day, month, year), goal.Message, goalStatus(goal.Achieved))
	}
	return table
}
