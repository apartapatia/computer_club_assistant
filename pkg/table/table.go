package table

import (
	"fmt"
	"time"
)

type Table struct {
	TableID    int
	ClientName string
	StartTime  time.Time
	AllTime    time.Duration
	Revenue    int
}

func NewTable(username string, tableID int, time time.Time) *Table {
	return &Table{
		TableID:    tableID,
		ClientName: username,
		StartTime:  time,
		AllTime:    0,
		Revenue:    0,
	}
}

func (t *Table) AllTimeFormatted() string {
	hours := int(t.AllTime.Hours())
	minutes := int(t.AllTime.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
