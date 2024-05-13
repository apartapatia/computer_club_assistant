package myparser

import (
	"os"
	"testing"
	"time"

	"github.com/apartapatia/computer_club_assistant/pkg/client"
	"github.com/apartapatia/computer_club_assistant/pkg/club"
)

func TestReadManagerEvents(t *testing.T) {
	file, err := os.CreateTemp("", "test_data")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(`10:00 1 anna
10:00 1 boris
10:00 2 anna 1
10:00 2 boris 2
11:00 1 charlie
11:00 2 charlie 3
12:00 1 david
12:00 1 emily
12:00 1 fiona
12:00 1 george
12:00 3 david
12:00 3 emily
12:00 3 fiona
12:00 3 george
14:00 4 anna
14:00 4 boris
14:00 4 charlie
16:00 4 david
16:00 2 emily 1
17:00 1 harry
17:00 2 harry 2
19:00 4 emily
19:00 4 fiona
22:00 1 kevin
22:00 2 kevin 1
`)
	if err != nil {
		t.Fatalf("Error writing test data to file: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("Error resetting file pointer: %v", err)
	}

	parser := NewFileParser(file)

	fakeClub := &club.Club{MaxTables: 4}

	expectedFirstClientTime, _ := time.Parse("15:04", "10:00")
	expectedLastClientTime, _ := time.Parse("15:04", "22:00")

	managers, err := parser.ReadManagerEvents(fakeClub)
	if err != nil {
		t.Fatalf("ReadManagerEvents returned error: %v", err)
	}

	expectedManagerCount := 25
	if len(managers) != expectedManagerCount {
		t.Errorf("Expected %d managers, got %d", expectedManagerCount, len(managers))
	}

	expectedFirstManager := &club.Manager{
		Time: expectedFirstClientTime,
		ID:   1,
		Client: &client.Client{
			Username: "anna",
			State:    1,
		},
	}

	if !managers[0].Time.Equal(expectedFirstManager.Time) ||
		managers[0].ID != expectedFirstManager.ID ||
		managers[0].Client.Username != expectedFirstManager.Client.Username ||
		managers[0].Client.State != expectedFirstManager.Client.State {
		t.Errorf("Expected first managerEvent to be %+v, got %+v", expectedFirstManager, managers[0])
	}

	expectedLastManager := &club.Manager{
		Time: expectedLastClientTime,
		ID:   2,
		Client: &client.Client{
			Username: "kevin",
			State:    2,
		},
		TableID: 1,
	}
	if !managers[24].Time.Equal(expectedLastManager.Time) ||
		managers[24].ID != expectedLastManager.ID ||
		managers[24].Client.Username != expectedLastManager.Client.Username ||
		managers[24].Client.State != expectedLastManager.Client.State ||
		managers[24].TableID != expectedLastManager.TableID {
		t.Errorf("Expected last managerEvent to be %+v, got %+v", expectedLastManager, managers[24])
	}
}

func TestReadClubInfo(t *testing.T) {
	file, err := os.CreateTemp("", "test_data")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(`3
10:00 23:00
1
`)
	if err != nil {
		t.Fatalf("Error writing test data to file: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("Error resetting file pointer: %v", err)
	}

	parser := NewFileParser(file)

	clubInfo, err := parser.ReadClubInfo()
	if err != nil {
		t.Fatalf("ReadClubInfo returned error: %v", err)
	}

	expectedMaxTables := 3
	if clubInfo.MaxTables != expectedMaxTables {
		t.Errorf("Expected MaxTables to be %d, got %d", expectedMaxTables, clubInfo.MaxTables)
	}

	expectedStartTime, _ := time.Parse("15:04", "10:00")
	expectedEndTime, _ := time.Parse("15:04", "23:00")

	if !clubInfo.WorkingTime.Open.Equal(expectedStartTime) || !clubInfo.WorkingTime.Close.Equal(expectedEndTime) {
		t.Errorf("Expected WorkingTime to be from %v to %v, got from %v to %v", expectedStartTime, expectedEndTime, clubInfo.WorkingTime.Open, clubInfo.WorkingTime.Close)
	}

	expectedPrice := 1
	if clubInfo.Price != expectedPrice {
		t.Errorf("Expected Price to be %d, got %d", expectedPrice, clubInfo.Price)
	}
}
