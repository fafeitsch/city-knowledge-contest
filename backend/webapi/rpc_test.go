package webapi

import (
	"encoding/json"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"github.com/stretchr/testify/assert"
	"testing"
)

var container = roomContainer{openRooms: make(map[string]contest.Room)}

func Test_createRoom(t *testing.T) {
	t.Run("validate too short user name", func(t *testing.T) {
		t.Parallel()
		request := createRoomRequest{Name: ""}
		js, _ := json.Marshal(request)
		_, err := container.createRoom(js)
		assert.EqualError(t, err, "a player name must not be empty")
	})
	t.Run("create and join room", func(t *testing.T) {
		request := createRoomRequest{Name: "JohnDoe"}
		js, _ := json.Marshal(request)
		processor, err := container.createRoom(js)
		assert.NoError(t, err)
		result, err := processor.process()
		assert.NoError(t, err)
		room := result.(createRoomResponse)
		assert.True(t, len(room.RoomKey) > 1, "room key \"%s\" too short", room.RoomKey)
		assert.True(t, len(room.PlayerKey) > 1, "player key \"%s\" too short", room.PlayerKey)
		assert.True(t, len(room.PlayerSecret) > 1, "player secret \"%s\" too short", room.PlayerSecret)
		assert.Equal(t, []string{"tooFewCoordinates"}, room.Errors)
		assert.Equal(t, 10, room.NumberOfQuestions)
		assert.Equal(t, 0, len(room.Area))
	})
}
