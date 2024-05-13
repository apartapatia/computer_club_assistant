package club

import (
	"fmt"
	"time"

	"github.com/apartapatia/computer_club_assistant/pkg/client"
)

type Manager struct {
	Time    time.Time
	ID      int
	Client  *client.Client
	TableID int
}

func (m *Manager) String() string {
	if m.TableID != 0 {
		return fmt.Sprintf("%s %d %s %d\n", m.Time.Format(TimeFormat), m.ID, m.Client.Username, m.TableID)
	} else {
		return fmt.Sprintf("%s %d %s\n", m.Time.Format(TimeFormat), m.ID, m.Client.Username)
	}
}

func (m *Manager) ErrorString(err error, ID int) string {
	return fmt.Sprintf("%s %d %s\n", m.Time.Format(TimeFormat), ID, err)
}

func (m *Manager) OutgoingString(idEvent int, username string, tableID int) string {
	if idEvent == 12 {
		return fmt.Sprintf("%s %d %s %d\n", m.Time.Format(TimeFormat), idEvent, username, tableID)
	} else {
		return fmt.Sprintf("%s %d %s\n", m.Time.Format(TimeFormat), idEvent, username)
	}
}

func NewManager(time time.Time, id int, clientName string, tableID int) *Manager {
	return &Manager{
		Time: time,
		ID:   id,
		Client: &client.Client{
			Username: clientName,
			State:    id,
		},
		TableID: tableID,
	}
}
