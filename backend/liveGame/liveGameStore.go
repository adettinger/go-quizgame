package livegame

import (
	"errors"
	"slices"
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

func (lgs *LiveGameStore) RemovePlayerByName(name string) error {
	lgs.mutex.Lock()
	defer lgs.mutex.Unlock()
	prevPlayerCount := len(lgs.Players)
	lgs.Players = slices.DeleteFunc(lgs.Players, func(p LivePlayer) bool {
		return p.Name == name
	})
	if prevPlayerCount == len(lgs.Players) {
		return errors.New("Cannot remove player: Player not found")
	}
	return nil
}

func (lgs *LiveGameStore) GetPlayerByName(name string) (LivePlayer, error) {
	lgs.mutex.RLock()
	defer lgs.mutex.RUnlock()

	for _, p := range lgs.Players {
		if p.Name == name {
			return p, nil
		}
	}
	return LivePlayer{}, errors.New("Player not found")
}


func (lgs *LiveGameStore) PlayerExistsByName(name string) bool {
	_, err := lgs.GetPlayerByName(name)
	return err == nil
}
