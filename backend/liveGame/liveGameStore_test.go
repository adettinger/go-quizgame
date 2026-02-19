package livegame_test

import (
	"slices"
	"testing"

	livegame "github.com/adettinger/go-quizgame/liveGame"
	"github.com/adettinger/go-quizgame/testutils"
	"github.com/google/uuid"
)

func TestAddPlayer(t *testing.T) {
	t.Run("Add multiple players", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		alexId, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		bobId, err := store.AddPlayer("Bob")
		testutils.AssertNoError(t, err)

		testutils.AssertTrue(t, store.PlayerExistsByName("Alex"))
		testutils.AssertTrue(t, store.PlayerExistsById(alexId))

		testutils.AssertTrue(t, store.PlayerExistsByName("Bob"))
		testutils.AssertTrue(t, store.PlayerExistsById(bobId))
	})

	t.Run("Add duplicate players", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		playerId, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		badId, err := store.AddPlayer("Alex")
		testutils.AssertHasError(t, err)
		testutils.AssertEqual(t, badId, uuid.Nil)

		testutils.AssertTrue(t, store.PlayerExistsByName("Alex"))
		testutils.AssertTrue(t, store.PlayerExistsById(playerId))
	})
}

func TestRemovePlayerByName(t *testing.T) {
	t.Run("Removes player if found", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		playerUUID, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		err = store.RemovePlayerByName("Alex")
		testutils.AssertNoError(t, err)

		testutils.AssertFalse(t, store.PlayerExistsByName("Alex"))
		testutils.AssertFalse(t, store.PlayerExistsById(playerUUID))
	})

	t.Run("Error if player not found", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		err := store.RemovePlayerByName("Alex")
		testutils.AssertHasError(t, err)

		testutils.AssertFalse(t, store.PlayerExistsByName("Alex"))
	})
}

func TestGetPlayerNameList(t *testing.T) {
	t.Run("Gets empty list for no names", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		got := store.GetPlayerNameList()
		testutils.AssertEqual(t, len(got), 0)
		testutils.AssertTrue(t, slices.Equal(got, []string{}))
	})
	t.Run("Gets name list", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		_, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		_, err = store.AddPlayer("Bob")
		testutils.AssertNoError(t, err)

		got := store.GetPlayerNameList()
		testutils.AssertEqual(t, len(got), 2)
		testutils.AssertTrue(t, slices.Equal(got, []string{"Alex", "Bob"}))
	})
}

func TestGetPlayerByName(t *testing.T) {
	t.Run("Gets player if exists", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		playerId, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		player, err := store.GetPlayerByName("Alex")
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, player.Id, playerId)
		testutils.AssertEqual(t, player.Name, "Alex")
	})

	t.Run("Error if player does not exist", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		player, err := store.GetPlayerByName("Alex")
		testutils.AssertHasError(t, err)

		testutils.AssertEqual(t, player, livegame.LivePlayer{})
	})
}

func TestPlayerExistsByName(t *testing.T) {
	t.Run("True if player exists", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		_, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		testutils.AssertTrue(t, store.PlayerExistsByName("Alex"))
	})

	t.Run("False if player does not exists", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		testutils.AssertFalse(t, store.PlayerExistsByName("Alex"))
	})
}

func TestGetPlayerById(t *testing.T) {
	t.Run("Gets player if exists", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		playerId, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		player, err := store.GetPlayerById(playerId)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, player.Id, playerId)
		testutils.AssertEqual(t, player.Name, "Alex")
	})

	t.Run("Error if player does not exist", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		player, err := store.GetPlayerById(uuid.New())
		testutils.AssertHasError(t, err)

		testutils.AssertEqual(t, player, livegame.LivePlayer{})
	})
}

func TestPlayerExistsById(t *testing.T) {
	t.Run("True if player exists", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		playerId, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		testutils.AssertTrue(t, store.PlayerExistsById(playerId))
	})

	t.Run("False if player does not exists", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		testutils.AssertFalse(t, store.PlayerExistsById(uuid.New()))
	})
}

func TestCreatePlayerId(t *testing.T) {
	t.Run("gets new unused id", func(t *testing.T) {
		store := livegame.NewLiveGameStore()

		_, err := store.AddPlayer("Alex")
		testutils.AssertNoError(t, err)

		id := store.CreatePlayerId()
		testutils.AssertFalse(t, store.PlayerExistsById(id))
	})
}
