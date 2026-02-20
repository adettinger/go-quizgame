package livegame

import (
	"errors"
	"slices"
	"sync"

	"github.com/adettinger/go-quizgame/types"
	"github.com/google/uuid"
)

type LivePlayer struct {
	Id   uuid.UUID
	Name string
}

type LiveGameStore struct {
	players []LivePlayer
	mutex   sync.RWMutex
}

func NewLiveGameStore() *LiveGameStore {
	return &LiveGameStore{}
}

func (lgs *LiveGameStore) AddPlayer(name string) (uuid.UUID, error) {
	if lgs.PlayerExistsByName(name) {
		return uuid.Nil, &types.ErrDuplicatePlayerName{PlayerName: name}
	}
	newPlayer := LivePlayer{
		Id:   lgs.CreatePlayerId(),
		Name: name,
	}
	lgs.mutex.Lock()
	defer lgs.mutex.Unlock()
	lgs.players = append(lgs.players, newPlayer)
	return newPlayer.Id, nil
}

func (lgs *LiveGameStore) RemovePlayerByName(name string) error {
	lgs.mutex.Lock()
	defer lgs.mutex.Unlock()
	prevPlayerCount := len(lgs.players)
	lgs.players = slices.DeleteFunc(lgs.players, func(p LivePlayer) bool {
		return p.Name == name
	})
	if prevPlayerCount == len(lgs.players) {
		return errors.New("Cannot remove player: Player not found")
	}
	return nil
}

func (lgs *LiveGameStore) GetPlayerNameList() []string {
	lgs.mutex.RLock()
	defer lgs.mutex.RUnlock()

	playerList := make([]string, len(lgs.players))
	for i, p := range lgs.players {
		playerList[i] = p.Name
	}
	return playerList
}

func (lgs *LiveGameStore) GetPlayerByName(name string) (LivePlayer, error) {
	lgs.mutex.RLock()
	defer lgs.mutex.RUnlock()

	for _, p := range lgs.players {
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

func (lgs *LiveGameStore) GetPlayerById(id uuid.UUID) (LivePlayer, error) {
	lgs.mutex.RLock()
	defer lgs.mutex.RUnlock()

	for _, p := range lgs.players {
		if p.Id == id {
			return p, nil
		}
	}
	return LivePlayer{}, errors.New("Player not found")
}

func (lgs *LiveGameStore) PlayerExistsById(id uuid.UUID) bool {
	_, err := lgs.GetPlayerById(id)
	return err == nil
}

func (lsg *LiveGameStore) CreatePlayerId() uuid.UUID {
	for {
		id := uuid.New()
		if !lsg.PlayerExistsById(id) {
			return id
		}
	}
}
