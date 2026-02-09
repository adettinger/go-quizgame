package webserver

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SessionStore struct {
	Sessions []Session
	mu       sync.RWMutex
}

type Session struct {
	Id      uuid.UUID
	Timeout time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		Sessions: make([]Session, 0),
		mu:       sync.RWMutex{},
	}
}

func (ss *SessionStore) CreateSession(duration time.Duration) Session {
	session := Session{ss.getNewId(), time.Now().Add(duration)}
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	ss.Sessions = append(ss.Sessions, session)
	return session
}

func (ss *SessionStore) GetBySessionId(id uuid.UUID) (Session, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	for _, s := range ss.Sessions {
		if s.Id == id {
			return s, nil
		}
	}
	return Session{}, errors.New("Problem not found")
}

func (ss *SessionStore) getNewId() uuid.UUID {
	for {
		uuid := uuid.New()
		if !ss.sessionExists(uuid) {
			return uuid
		}
	}
}

func (ss *SessionStore) sessionExists(uuid uuid.UUID) bool {
	_, err := ss.GetBySessionId(uuid)
	return err == nil
}
