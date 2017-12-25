package cmd

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gosuri/uitable"
	_ "github.com/mattn/go-sqlite3"
)

const (
	CREATETABLE = `
    CREATE TABLE goal (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
        timestamp TEXT,
        goalmsg TEXT, 
        achieved INTEGER DEFAULT 0
        )
`

	INSERTGOAL    = `INSERT INTO goal (timestamp, goalmsg, achieved) VALUES (?, ?, ?)`
	GETGOALS      = `SELECT timestamp, goalmsg, achieved FROM goal`
	GETTODAYGOALS = "SELECT * FROM goal WHERE timestamp=?"
	DROPTABLE     = `DROP TABLE goal`
	UPDATEGOAL    = `UPDATE goal SET achieved=1 WHERE goalmsg=(?)`
)

type Goal struct {
	Message   string
	Timestamp time.Time
	Achieved  bool
}

type Goals []*Goal

var DBFILE = UserHomeDir() + "/.gftd.db"

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

// Create a database file with a table.
//
func CreateDB() error {
	db, err := sql.Open("sqlite3", DBFILE)
	if err != nil {
		return err
	}
	defer db.Close()

	statement, err := db.Prepare(CREATETABLE)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
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

func ReadTodayGoals() ([]*Goal, error) {
	db, err := sql.Open("sqlite3", DBFILE)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	y, m, d := time.Now().Date()

	rows, err := db.Query(GETTODAYGOALS, fmt.Sprintf("%v-%v-%v", y, m, d))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*Goal
	var isAchieved int
	var goalmsg, timestamp string

	for rows.Next() {
		err = rows.Scan(&timestamp, &goalmsg, &isAchieved)
		if err != nil {
			return nil, err
		}
		times, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timestamp)
		if err != nil {
			return nil, err
		}

		if isSameDate(times, time.Now()) {
			isachieved := isAchieved != 0
			goals = append(goals, &Goal{goalmsg, times, isachieved})
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return goals, nil
}

// Reads all the goals currently in the database
//
func ReadAllGoals() ([]*Goal, error) {
	db, err := sql.Open("sqlite3", DBFILE)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(GETGOALS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*Goal
	var isAchieved int
	var goalmsg, timestamp string

	for rows.Next() {
		err = rows.Scan(&timestamp, &goalmsg, &isAchieved)
		if err != nil {
			return nil, err
		}
		times, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, err
		}
		isachieved := isAchieved != 0
		g := Goal{goalmsg, times, isachieved}
		//fmt.Println("Reading goal:", g)
		goals = append(goals, &g)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return goals, nil
}

// Writes a goal to the database
//
func WriteGoal(goal *Goal) error {
	db, err := sql.Open("sqlite3", DBFILE)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(INSERTGOAL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	isAchieved := 0
	if goal.Achieved {
		isAchieved = 1
	}

	_, er := stmt.Exec(goal.Timestamp.Format(time.RFC3339), goal.Message, isAchieved)
	tx.Commit()

	return er
}

func UpdateGoal(msg string) error {
	db, err := sql.Open("sqlite3", DBFILE)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(UPDATEGOAL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, er := stmt.Exec(msg)
	tx.Commit()

	return er
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
	table.AddRow("----", strings.Repeat("-", 16), strings.Repeat("-", 50), "--------")
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
