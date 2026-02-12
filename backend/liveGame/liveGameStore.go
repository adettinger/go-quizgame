package livegame

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type LivePlayer struct {
	Id   uuid.UUID
	Name string
}

type LiveGameStore struct {
	Players []LivePlayer
	mutex   sync.RWMutex
}

func NewLiveGameStore() *LiveGameStore {
	return &LiveGameStore{}
}

func (lgs *LiveGameStore) AddPlayer(name string) (uuid.UUID, error) {
	if lgs.PlayerExistsByName(name) {
		return uuid.Nil, errors.New("Player with name already exists")
	}
	newPlayer := LivePlayer{
		Id:   uuid.New(),
		Name: name,
	}
	lgs.mutex.Lock()
	defer lgs.mutex.Unlock()
	lgs.Players = append(lgs.Players, newPlayer)
	return newPlayer.Id, nil
}

func (lgs *LiveGameStore) PlayerExistsByName(name string) bool {
	lgs.mutex.RLock()
	defer lgs.mutex.RUnlock()

	for _, p := range lgs.Players {
		if p.Name == name {
			return true
		}
	}
	return false
}
