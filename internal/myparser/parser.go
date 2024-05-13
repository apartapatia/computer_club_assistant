package myparser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apartapatia/computer_club_assistant/pkg/client"
	"github.com/apartapatia/computer_club_assistant/pkg/club"
)

var (
	ErrParseInt = errors.New("ParseIntError")
	ErrReadData = errors.New("ReadDataError")
)

type Parser interface {
	ReadClubInfo() (*club.Club, error)
	ReadManagerEvents(activeClub *club.Club) ([]*club.Manager, error)
	ParseInt(data string) (int, error)
	InvalidParse(line []string)
}

type FileParser struct {
	scanner *bufio.Scanner
}

func (fp *FileParser) ParseInt(data string) (int, error) {
	v, err := strconv.Atoi(data)
	if err != nil || v <= 0 {
		return 0, ErrParseInt
	}
	return v, nil
}

func (fp *FileParser) InvalidParse(line []string) {
	fmt.Println(line)
	os.Exit(1)
}

func (fp *FileParser) ReadClubInfo() (*club.Club, error) {
	if !fp.scanner.Scan() {
		return nil, ErrReadData
	}

	maxTablesData := fp.scanner.Text()
	if len(maxTablesData) == 0 {
		fp.InvalidParse([]string{maxTablesData})
	}

	maxTables, err := fp.ParseInt(maxTablesData)
	if err != nil {
		fp.InvalidParse([]string{maxTablesData})
	}

	if !fp.scanner.Scan() {
		return nil, ErrReadData
	}

	workingTimeData := fp.scanner.Text()
	if len(workingTimeData) == 0 {
		fp.InvalidParse([]string{workingTimeData})
	}

	times := strings.Split(workingTimeData, " ")
	if len(times) != 2 {
		fp.InvalidParse(times)
	}

	startTime, err := time.Parse(club.TimeFormat, times[0])
	if err != nil {
		fp.InvalidParse(times)
	}

	endTime, err := time.Parse(club.TimeFormat, times[1])
	if err != nil {
		fp.InvalidParse(times)
	}

	if endTime.Before(startTime) {
		fp.InvalidParse(times)
	}
	workingTime := club.NewWorkingTime(startTime, endTime)

	if !fp.scanner.Scan() {
		return nil, ErrReadData
	}

	priceData := fp.scanner.Text()
	if len(priceData) == 0 {
		fp.InvalidParse([]string{priceData})
	}

	price, err := fp.ParseInt(priceData)
	if err != nil {
		fp.InvalidParse([]string{priceData})
	}

	return club.NewClub(workingTime, price, maxTables), nil
}

func (fp *FileParser) ReadManagerEvents(activeClub *club.Club) ([]*club.Manager, error) {
	var (
		managers []*club.Manager
	)

	for fp.scanner.Scan() {
		line := fp.scanner.Text()

		parts := strings.Fields(line)
		if len(parts) < 3 || len(parts) > 4 {
			fp.InvalidParse([]string{line})
		}

		eventTime, err := time.Parse(club.TimeFormat, parts[0])
		if err != nil {
			fp.InvalidParse([]string{line})
		}

		eventType, err := fp.ParseInt(parts[1])
		if err != nil || eventType == 2 && len(parts) < 4 {
			fp.InvalidParse([]string{line})
		}

		if ok, _ := client.ValidateUsername(parts[2]); !ok {
			fp.InvalidParse([]string{line})
		}

		clientName := parts[2]

		var tableID int
		if len(parts) > 3 {
			if eventType != 2 {
				fp.InvalidParse([]string{line})
			}

			tableID, err = fp.ParseInt(parts[3])
			if err != nil || tableID > activeClub.MaxTables {
				fp.InvalidParse([]string{line})
			}
		}

		manager := club.NewManager(eventTime, eventType, clientName, tableID)
		managers = append(managers, manager)
	}

	if err := fp.scanner.Err(); err != nil {
		return nil, err
	}

	sort.Slice(managers, func(i, j int) bool {
		return managers[i].Time.Before(managers[j].Time)
	})

	return managers, nil
}

func NewFileParser(file *os.File) *FileParser {
	return &FileParser{
		scanner: bufio.NewScanner(file),
	}
}
