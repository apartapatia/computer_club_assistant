package client

import (
	"errors"
	"sort"
	"sync"
)

const QueueState = 3

var (
	ErrClientAlreadyExists = errors.New("YouShallNotPass")
	ErrClientNotFound      = errors.New("ClientUnknown")
	ErrClientIsNil         = errors.New("ClientIsNil")
)

type ClientRepository interface {
	Add(client *Client) error
	Exists(username string) bool
	Get(username string) (*Client, error)
	Remove(username string) error
	UpdateStatus(username string, newStatus int) error
	Queue() []*Client
	GetAll() map[string]*Client
}

type ClientRepositoryMemory struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

func (cr *ClientRepositoryMemory) GetAll() map[string]*Client {
	return cr.clients
}

func (cr *ClientRepositoryMemory) Add(client *Client) error {
	if client == nil {
		return ErrClientIsNil
	}

	cr.mu.Lock()
	defer cr.mu.Unlock()

	if _, ok := cr.clients[client.Username]; ok {
		return ErrClientAlreadyExists
	}

	cr.clients[client.Username] = client
	return nil
}

func (cr *ClientRepositoryMemory) Exists(username string) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	_, exists := cr.clients[username]
	return exists
}

func (cr *ClientRepositoryMemory) Get(username string) (*Client, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	client, ok := cr.clients[username]
	if !ok {
		return &Client{}, ErrClientNotFound
	}
	return client, nil
}

func (cr *ClientRepositoryMemory) Queue() []*Client {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var keys []string
	for key := range cr.clients {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var queueClients []*Client

	for _, key := range keys {
		client := cr.clients[key]
		if client.State == QueueState {
			queueClients = append(queueClients, client)
		}
	}

	return queueClients
}

func (cr *ClientRepositoryMemory) Remove(username string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if _, ok := cr.clients[username]; !ok {
		return ErrClientNotFound
	}
	delete(cr.clients, username)
	return nil
}

func (cr *ClientRepositoryMemory) UpdateStatus(username string, newStatus int) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	client, ok := cr.clients[username]
	if !ok {
		return ErrClientNotFound
	}
	client.State = newStatus
	return nil
}

func NewMemoryRepo() *ClientRepositoryMemory {
	return &ClientRepositoryMemory{
		clients: make(map[string]*Client),
	}
}
