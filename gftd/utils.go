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
)

type Goal struct {
	Message   string    `json:"Message"`
	Timestamp time.Time `json:"Timestamp"`
	Achieved  bool      `json:"Achieved"`
}

var DBFILE = UserHomeDir() + "/.gftd.json"

func init() {
	if exists, err := IsDBExists(); exists {
		if err == nil {
			return
		}
	}

	err := CreateDB()
	if err != nil {
		fmt.Println("Error while creating the database:", err)
	}
}

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
func ReadAllGoals(file io.Reader) ([]Goal, error) {
	var goals []Goal
	dec := json.NewDecoder(file)
	if err := dec.Decode(&goals); err != nil && err != io.EOF {
		return nil, err // TODO: Handle the case of empty file EOF error
	}

	for i, goal := range goals {
		fmt.Printf("%d. %v", i, goal)
	}

	return goals, nil
}

// Writes a goal to the database
//
func WriteAllGoals(file io.Writer, goals []Goal) error {
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
