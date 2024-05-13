package handlers

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/apartapatia/computer_club_assistant/pkg/client"
	"github.com/apartapatia/computer_club_assistant/pkg/club"
	"github.com/apartapatia/computer_club_assistant/pkg/table"
)

var (
	ErrClientIsWaiting = errors.New("ICanWaitNoLonger!")
	ErrNotOpen         = errors.New("NotOpenYet")
)

const (
	IncomingClientCome         = 1
	IncomingClientTookTheTable = 2
	IncomingClientIsWaiting    = 3
	IncomingClientLeft         = 4

	OutgoingClientAfterClose               = 11
	OutgoingClientTokeTheTableAfterWaiting = 12
	OutgoingClientError                    = 13
)

type CommandHandler struct {
	Club     *club.Club
	Managers []*club.Manager
	Clients  client.ClientRepository
	Tables   table.TableRepository
}

func (h *CommandHandler) HandleCommands() string {
	var sb strings.Builder
	sb.WriteString(h.Club.WorkingTime.Open.Format(club.TimeFormat) + "\n")

	for _, m := range h.Managers {
		if !club.IsTimeWithinWorkingHours(*h.Club.WorkingTime, m.Time) {
			sb.WriteString(m.String())
			sb.WriteString(m.ErrorString(ErrNotOpen, OutgoingClientError))
			continue
		}

		switch m.ID {
		case IncomingClientCome:
			sb.WriteString(h.handleIncomingClientCome(m))
		case IncomingClientTookTheTable:
			sb.WriteString(h.handleIncomingClientTookTheTable(m))
		case IncomingClientIsWaiting:
			sb.WriteString(h.handleIncomingClientIsWaiting(m))
		case IncomingClientLeft:
			sb.WriteString(h.handleIncomingClientLeft(m))
		}
	}

	sb.WriteString(h.checkLastClient())
	sb.WriteString(h.Club.WorkingTime.Close.Format(club.TimeFormat) + "\n")
	sb.WriteString(h.calculateRevenue())
	return sb.String()[:sb.Len()-1]
}

func (h *CommandHandler) handleIncomingClientCome(manager *club.Manager) string {
	var sb strings.Builder
	sb.WriteString(manager.String())

	if err := h.Clients.Add(manager.Client); err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
	}
	return sb.String()
}

func (h *CommandHandler) handleIncomingClientTookTheTable(manager *club.Manager) string {
	var sb strings.Builder
	sb.WriteString(manager.String())

	c, err := h.Clients.Get(manager.Client.Username)
	if err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
		return sb.String()
	}

	if err := h.Clients.UpdateStatus(c.Username, manager.ID); err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
	}

	if currentID, ok := h.Tables.Exists(c.Username); ok {
		h.Tables.TakeDownTable(c.Username)
		err = h.Tables.UpdateRevenue(currentID, h.Club.Price, manager.Time)
		if err != nil {
			return err.Error()
		}
	}

	if err := h.Tables.TakeUpTable(c.Username, manager.TableID, manager.Time); err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
	}

	return sb.String()
}

func (h *CommandHandler) handleIncomingClientIsWaiting(manager *club.Manager) string {
	var sb strings.Builder
	sb.WriteString(manager.String())

	c, err := h.Clients.Get(manager.Client.Username)
	if err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
		return sb.String()
	}

	if err := h.Clients.UpdateStatus(c.Username, manager.ID); err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
	}

	if h.Tables.CountEmptyTables() != 0 {
		sb.WriteString(manager.ErrorString(ErrClientIsWaiting, OutgoingClientError))
	}

	if len(h.Clients.Queue()) > h.Club.MaxTables {
		if err := h.Clients.Remove(c.Username); err != nil {
			sb.WriteString(manager.ErrorString(err, OutgoingClientError))
		}
		sb.WriteString(manager.OutgoingString(OutgoingClientAfterClose, c.Username, manager.TableID))
	}

	return sb.String()
}

func (h *CommandHandler) handleIncomingClientLeft(manager *club.Manager) string {
	var sb strings.Builder
	sb.WriteString(manager.String())

	c, err := h.Clients.Get(manager.Client.Username)
	if err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
		return sb.String()
	}

	tableID := h.Tables.TakeDownTable(c.Username)

	if err := h.Clients.Remove(c.Username); err != nil {
		sb.WriteString(manager.ErrorString(err, OutgoingClientError))
	}

	if tableID != 0 {
		err = h.Tables.UpdateRevenue(tableID, h.Club.Price, manager.Time)
		if err != nil {
			return err.Error()
		}
	}

	if len(h.Clients.Queue()) != 0 && tableID != 0 {
		usernameFirstQueue := h.Clients.Queue()[0].Username

		if err := h.Tables.TakeUpTable(usernameFirstQueue, tableID, manager.Time); err != nil {
			sb.WriteString(manager.ErrorString(err, OutgoingClientError))
		}

		if err := h.Clients.UpdateStatus(usernameFirstQueue, IncomingClientTookTheTable); err != nil {
			sb.WriteString(manager.ErrorString(err, OutgoingClientError))
		}

		sb.WriteString(manager.OutgoingString(OutgoingClientTokeTheTableAfterWaiting, usernameFirstQueue, tableID))
	}

	return sb.String()
}

func (h *CommandHandler) calculateRevenue() string {
	var sb strings.Builder
	tables := h.Tables.GetAll()

	for tableID := 1; tableID <= h.Club.MaxTables; tableID++ {
		if t, ok := tables[tableID]; ok {
			sb.WriteString(fmt.Sprintf("%d %d %s\n", t.TableID, t.Revenue, t.AllTimeFormatted()))
		} else {
			sb.WriteString(fmt.Sprintf("%d %d %s\n", tableID, 0, "00:00"))
		}
	}

	return sb.String()
}

func (h *CommandHandler) checkLastClient() string {
	var sb strings.Builder
	clients := h.Clients.GetAll()

	var queueClientNames []string
	for _, c := range clients {
		if c.Username != "" {
			queueClientNames = append(queueClientNames, c.Username)
		}
	}

	sort.Strings(queueClientNames)

	for _, clientName := range queueClientNames {
		sb.WriteString(fmt.Sprintf("%s %d %s\n", h.Club.WorkingTime.Close.Format(club.TimeFormat), OutgoingClientAfterClose, clientName))

		if currentID, ok := h.Tables.Exists(clientName); ok {
			err := h.Tables.UpdateRevenue(currentID, h.Club.Price, h.Club.WorkingTime.Close)
			if err != nil {
				return err.Error()
			}

			h.Tables.TakeDownTable(clientName)
		}

		err := h.Clients.Remove(clientName)
		if err != nil {
			return err.Error()
		}
	}

	return sb.String()
}

func NewCommandHandler(club *club.Club, managers []*club.Manager, clients client.ClientRepository, tables table.TableRepository) *CommandHandler {
	return &CommandHandler{
		Club:     club,
		Managers: managers,
		Clients:  clients,
		Tables:   tables,
	}
}
