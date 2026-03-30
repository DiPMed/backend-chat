package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{sessions: make(map[string]*Session)}
}

func (m *MemoryStore) Create() (*Session, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}

	s := &Session{ID: id}

	m.mu.Lock()
	m.sessions[id] = s
	m.mu.Unlock()

	return s, nil
}

func (m *MemoryStore) Get(id string) (*Session, error) {
	m.mu.RLock()
	s, ok := m.sessions[id]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("session %q not found", id)
	}
	return s, nil
}

func (m *MemoryStore) Append(id string, msgs ...Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[id]
	if !ok {
		return fmt.Errorf("session %q not found", id)
	}

	s.Messages = append(s.Messages, msgs...)
	return nil
}

func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
