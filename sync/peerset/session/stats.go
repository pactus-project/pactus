package session

import "fmt"

type Stats struct {
	Total       int
	Open        int
	Completed   int
	Uncompleted int
}

func (ss *Stats) String() string {
	return fmt.Sprintf("total: %v, open: %v, completed: %v, uncompleted: %v",
		ss.Total, ss.Open, ss.Completed, ss.Uncompleted)
}
