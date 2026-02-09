package webserver

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SessionStore struct {
	Sessions map[uuid.UUID]SessionData
	mu       sync.RWMutex
}

type SessionData struct {
	Timeout time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		Sessions: make(map[uuid.UUID]SessionData, 0),
		mu:       sync.RWMutex{},
	}
}

func (ss *SessionStore) CreateSession(duration time.Duration) (uuid.UUID, SessionData) {
	sessionID := ss.getNewId()
	sessionData := SessionData{time.Now().Add(duration)}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.Sessions[sessionID] = sessionData
	return sessionID, sessionData
}

func (ss *SessionStore) GetBySessionId(id uuid.UUID) (SessionData, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	s, ok := ss.Sessions[id]
	if !ok {
		return SessionData{}, errors.New("Session not found")
	}
	return s, nil
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

func (ss *SessionStore) IsSessionActive(id uuid.UUID, time time.Time) (bool, error) {
	session, err := ss.GetBySessionId(id)
	if err != nil {
		return false, errors.New("Session not found")
	}
	return session.Timeout.After(time), nil
}

func (ss *SessionStore) DeleteSession(id uuid.UUID) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	delete(ss.Sessions, id)
}
