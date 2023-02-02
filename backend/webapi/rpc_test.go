package webapi

import (
	"encoding/json"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var container = roomContainer{openRooms: make(map[string]contest.Room)}

func TestMain(m *testing.M) {
	osrmServer := createOsrmTestServer()
	defer osrmServer.Close()
	contest.OsrmServer = osrmServer.URL
	m.Run()
}

func Test_createRoom(t *testing.T) {
	t.Parallel()
	t.Run(
		"validate too short user name", func(t *testing.T) {
			t.Parallel()
			request := createRoomRequest{Name: ""}
			_, err := container.createRoom(mustMarshal(request))
			assert.EqualError(t, err, "a player name must not be empty")
		},
	)
	t.Run(
		"create and join room", func(t *testing.T) {
			t.Parallel()
			request := createRoomRequest{Name: "JohnDoe"}
			processor, err := container.createRoom(mustMarshal(request))
			assert.NoError(t, err)
			result, err := processor.process()
			assert.NoError(t, err)
			room := result.(createRoomResponse)
			assert.True(
				t, len(room.RoomKey) > 1, "room key \"%s\" too short", room.RoomKey,
			)
			assert.True(
				t, len(room.PlayerKey) > 1, "player key \"%s\" too short", room.PlayerKey,
			)
			assert.True(
				t, len(room.PlayerSecret) > 1, "player secret \"%s\" too short", room.PlayerSecret,
			)
			assert.Equal(t, []string{"tooFewCoordinates"}, room.Errors)
			assert.Equal(t, 10, room.NumberOfQuestions)
			assert.Equal(t, 0, len(room.Area))
		},
	)
}

func Test_updateRoom(t *testing.T) {
	t.Parallel()
	t.Run(
		"validate room not found", func(t *testing.T) {
			t.Parallel()
			request := roomUpdateRequest{RoomKey: "not found"}
			_, err := container.updateRoom(mustMarshal(request))
			assert.EqualError(t, err, "room with key \"not found\" not found")
		},
	)
	t.Run(
		"update room", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			player := room.Join("playerA")
			request := roomUpdateRequest{
				PlayerKey:         player.Key,
				PlayerSecret:      player.Secret,
				NumberOfQuestions: 42,
				RoomKey:           room.Key(),
				Area:              [][2]float64{{1, 2}, {3, 4}},
			}
			processor, err := container.updateRoom(mustMarshal(request))
			assert.NoError(t, err)
			result, err := processor.process()
			defer processor.release()
			assert.NoError(t, err)
			response := result.(updateRoomResponse)
			assert.Equal(t, []string{"tooFewCoordinates"}, response.Errors)
			got := room.Options()
			assert.Equal(t, contest.Coordinate{Lat: 1, Lng: 2}, got.Area[0])
			assert.Equal(t, 2, len(got.Area))
			assert.Equal(t, 42, got.NumberOfQuestions)
		},
	)
}

func Test_JoinRoom(t *testing.T) {
	t.Parallel()
	t.Run(
		"validate room not found", func(t *testing.T) {
			t.Parallel()
			request := joinRequest{RoomKey: "not found"}
			_, err := container.joinRoom(mustMarshal(request))
			assert.EqualError(t, err, "room with key \"not found\" not found")
		},
	)
	t.Run(
		"join room", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			request := joinRequest{Name: "Doc", RoomKey: room.Key()}
			processor, err := container.joinRoom(mustMarshal(request))
			assert.NoError(t, err)
			defer processor.release()
			result, err := processor.process()
			assert.NoError(t, err)
			got := result.(joinResponse)
			assert.True(t, len(got.PlayerKey) > 1)
			assert.True(t, len(got.PlayerSecret) > 1)
			assert.Equal(t, "Doc", got.Name)
			player, _ := room.FindPlayer(got.PlayerKey)
			assert.Equal(t, "Doc", player.Name)
		},
	)
}

func Test_StartGame(t *testing.T) {
	t.Parallel()
	t.Run(
		"validate player not joined", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    "xyz",
				PlayerSecret: "abc",
			}
			context, err := container.startGame(mustMarshal(request))
			assert.EqualError(
				t, err, "player with key \"xyz\" has not joined the room yet or secret is wrong",
			)
			assert.NotNil(t, context)
		},
	)
	t.Run(
		"room is not yet complete", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			player := room.Join("Doc")
			room.SetOptions(
				contest.RoomOptions{
					NumberOfQuestions: 1,
					Area:              []contest.Coordinate{},
				}, player.Key,
			)
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    player.Key,
				PlayerSecret: player.Secret,
			}
			_, err := container.startGame(mustMarshal(request))
			assert.EqualError(
				t, err, "the room is not yet configured correctly, "+
					"it still has the following config errors: [tooFewCoordinates]",
			)
		},
	)
	t.Run(
		"start a game", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			player := room.Join("Doc")
			room.SetOptions(
				contest.RoomOptions{
					NumberOfQuestions: 1,
					Area: []contest.Coordinate{
						{Lat: 1, Lng: 2},
						{Lat: 3, Lng: 4},
						{Lat: 1, Lng: 5},
					},
				}, player.Key,
			)
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    player.Key,
				PlayerSecret: player.Secret,
			}
			process, err := container.startGame(mustMarshal(request))
			assert.NoError(t, err)
			result, err := process.process()
			process.release()
			assert.NoError(t, err)
			assert.Equal(t, map[string]any{}, result)
			var group sync.WaitGroup
			group.Add(1)
			go waitUntilRoomHasActiveQuestion(room, &group)
			group.Wait()
			room.Lock()
			assert.True(t, room.HasActiveQuestion())
		},
	)
}

func Test_AnswerQuestion(t *testing.T) {
	t.Parallel()
	t.Run(
		"validate player not joined", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    "xyz",
				PlayerSecret: "abc",
			}
			context, err := container.answerQuestion(mustMarshal(request))
			assert.EqualError(
				t, err, "player with key \"xyz\" has not joined the room yet or secret is wrong",
			)
			assert.NotNil(t, context.release)
		},
	)
	t.Run(
		"room has no active question", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			player := room.Join("Doc")
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    player.Key,
				PlayerSecret: player.Secret,
			}
			context, err := container.answerQuestion(mustMarshal(request))
			assert.EqualError(
				t, err, "question cannot be answered because there is no active question",
			)
			assert.NotNil(t, context.release)
		},
	)
	t.Run(
		"answer a question to running game", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			player := room.Join("Doc")
			room.SetOptions(
				contest.RoomOptions{
					NumberOfQuestions: 1,
					Area:              []contest.Coordinate{{Lat: 1, Lng: 2}, {Lat: 3, Lng: 4}, {Lat: 1, Lng: 5}},
				}, player.Key,
			)
			room.Play(player.Key)
			var group sync.WaitGroup
			group.Add(1)
			go waitUntilRoomHasActiveQuestion(room, &group)
			group.Wait()
			request := guessRequest{
				RoomKey:      room.Key(),
				PlayerKey:    player.Key,
				PlayerSecret: player.Secret,
				Guess:        [2]float64{0.5, 0.3},
			}
			process, err := container.answerQuestion(mustMarshal(request))
			assert.NoError(t, err)
			result, err := process.process()
			points := result.(map[string]int)
			process.release()
			assert.NoError(t, err)
			assert.True(t, points["points"] > 0)
		},
	)
}

func Test_AdvanceGame(t *testing.T) {
	t.Parallel()
	t.Run(
		"validate player not joined", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    "xyz",
				PlayerSecret: "abc",
			}
			context, err := container.advanceGame(mustMarshal(request))
			assert.EqualError(
				t, err, "player with key \"xyz\" has not joined the room yet or secret is wrong",
			)
			assert.NotNil(t, context.release)
		},
	)
	t.Run(
		"room cannot be advanced", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			player := room.Join("Doc")
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    player.Key,
				PlayerSecret: player.Secret,
			}
			context, err := container.advanceGame(mustMarshal(request))
			assert.EqualError(t, err, "the room \""+room.Key()+"\" cannot be advanced")
			assert.NotNil(t, context.release)
		},
	)
	t.Run(
		"advance a game", func(t *testing.T) {
			t.Parallel()
			room := contest.NewRoom()
			container.Lock()
			container.openRooms[room.Key()] = room
			container.Unlock()
			player := room.Join("Doc")
			room.SetOptions(
				contest.RoomOptions{
					NumberOfQuestions: 3,
					Area:              []contest.Coordinate{{Lat: 1, Lng: 2}, {Lat: 3, Lng: 4}, {Lat: 1, Lng: 5}},
				}, player.Key,
			)
			room.Play(player.Key)
			var group sync.WaitGroup
			group.Add(1)
			go waitUntilRoomHasActiveQuestion(room, &group)
			group.Wait()
			_, err := room.AnswerQuestion(player.Key, contest.Coordinate{Lat: 0.5, Lng: 0.3})
			assert.NoError(t, err)
			room.Lock()
			assert.False(t, room.HasActiveQuestion())
			room.Unlock()
			request := startGameRequest{
				RoomKey:      room.Key(),
				PlayerKey:    player.Key,
				PlayerSecret: player.Secret,
			}
			process, err := container.advanceGame(mustMarshal(request))
			assert.NoError(t, err)
			_, err = process.process()
			assert.NoError(t, err)
			process.release()
			group = sync.WaitGroup{}
			group.Add(1)
			go waitUntilRoomHasActiveQuestion(room, &group)
			group.Wait()
		},
	)
}

func waitUntilRoomHasActiveQuestion(room contest.Room, group *sync.WaitGroup) {
	activeQuestion := false
	for !activeQuestion {
		room.Lock()
		activeQuestion = room.HasActiveQuestion()
		room.Unlock()
	}
	group.Done()
}
func createOsrmTestServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				response := mustMarshal(
					map[string]any{
						"waypoints": []any{
							map[string]any{"location": [2]float64{2, 3}, "name": "Somestreet"},
						},
					},
				)
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write(response)
			},
		),
	)
}

func mustMarshal(payload any) []byte {
	result, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return result
}
