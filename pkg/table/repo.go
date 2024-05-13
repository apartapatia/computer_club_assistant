package table

import (
	"errors"
	"math"
	"sync"
	"time"
)

var (
	ErrPlaceIsBusy   = errors.New("PlaceIsBusy")
	ErrTableNotFound = errors.New("TableNotFound")
	ErrTablesFull    = errors.New("TablesFull")
)

type TableRepository interface {
	GetAll() map[int]*Table
	TakeUpTable(clientName string, tableID int, t time.Time) error
	TakeDownTable(clientName string) int
	UpdateRevenue(tableID, price int, t time.Time) error
	Exists(clientName string) (int, bool)
	CountEmptyTables() int
}

type TableRepositoryMemory struct {
	tables    map[int]*Table
	maxTables int
	mu        *sync.RWMutex
}

func (r *TableRepositoryMemory) GetAll() map[int]*Table {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tables
}

func (r *TableRepositoryMemory) TakeUpTable(clientName string, tableID int, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tableID > r.maxTables {
		return ErrTablesFull
	}

	table, ok := r.tables[tableID]
	if !ok {
		table = NewTable("", tableID, t)
		r.tables[tableID] = table
	}

	if table.ClientName != "" {
		return ErrPlaceIsBusy
	}

	table.ClientName = clientName
	table.StartTime = t
	return nil
}

func (r *TableRepositoryMemory) TakeDownTable(clientName string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	for tableID, table := range r.tables {
		if table.ClientName == clientName {
			r.tables[tableID].ClientName = ""
			return tableID
		}
	}

	return 0
}

func (r *TableRepositoryMemory) UpdateRevenue(tableID, price int, t time.Time) error {
	table, ok := r.tables[tableID]
	if !ok {
		return ErrTableNotFound
	}

	duration := t.Sub(table.StartTime)
	durationInHours := duration.Minutes() / 60
	priceCounter := int(math.Ceil(durationInHours))

	table.Revenue += priceCounter * price
	table.AllTime += duration

	return nil
}

func (r *TableRepositoryMemory) Exists(clientName string) (int, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, table := range r.tables {
		if table.ClientName == clientName {
			return table.TableID, true
		}
	}
	return 0, false
}

func (r *TableRepositoryMemory) CountEmptyTables() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, table := range r.tables {
		if table.ClientName == "" {
			count++
		}
	}

	if len(r.tables) != r.maxTables {
		return r.maxTables
	}

	return count
}

func NewMemoryRepo(maxTables int) *TableRepositoryMemory {
	return &TableRepositoryMemory{
		tables:    make(map[int]*Table),
		maxTables: maxTables,
		mu:        &sync.RWMutex{},
	}
}
