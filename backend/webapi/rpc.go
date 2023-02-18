package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"log"
	"sync"
	"time"
)

type roomContainer struct {
	sync.RWMutex
	openRooms map[string]contest.Room
}

type rpcHandler func(message json.RawMessage) (*rpcRequestContext, error)

type rpcRequestContext struct {
	process func() (any, error)
	release func()
}

type createRoomRequest struct {
	Name string `json:"name"`
}

type createRoomResponse struct {
	RoomKey           string       `json:"roomKey"`
	PlayerKey         string       `json:"playerKey"`
	PlayerSecret      string       `json:"playerSecret"`
	NumberOfQuestions int          `json:"numberOfQuestions"`
	Area              [][2]float64 `json:"area"`
	Errors            []string     `json:"errors"`
}

func (r *roomContainer) createRoom(message json.RawMessage) (*rpcRequestContext, error) {
	request := parseMessage[createRoomRequest](message)
	if len(request.Name) == 0 {
		return nil, fmt.Errorf("a player name must not be empty")
	}
	return &rpcRequestContext{
		process: func() (any, error) {
			room := contest.NewRoom()
			player := room.Join(request.Name)
			r.Lock()
			r.openRooms[room.Key()] = room
			r.Unlock()
			area := make([][2]float64, 0, len(room.Options().Area))
			for _, coordinate := range room.Options().Area {
				area = append(area, [2]float64{coordinate.Lat, coordinate.Lng})
			}
			response := createRoomResponse{
				RoomKey:           room.Key(),
				PlayerKey:         player.Key,
				PlayerSecret:      player.Secret,
				Errors:            room.ConfigErrors(),
				NumberOfQuestions: room.Options().NumberOfQuestions,
				Area:              area,
			}
			log.Printf("Created room \"%s\" from player \"%s\" (\"%s\").", room.Key(), player.Secret, player.Name)
			return response, nil
		},
	}, nil
}

type roomUpdateRequest struct {
	Area              [][2]float64 `json:"area"`
	NumberOfQuestions int          `json:"numberOfQuestions"`
	RoomKey           string       `json:"roomKey"`
	MaxAnswerTimeSec  int          `json:"maxAnswerTimeSec"`
	PlayerKey         string       `json:"playerKey"`
	PlayerSecret      string       `json:"playerSecret"`
}

type updateRoomResponse struct {
	Errors []string `json:"errors"`
}

func (r *roomContainer) updateRoom(message json.RawMessage) (*rpcRequestContext, error) {
	request := parseMessage[roomUpdateRequest](message)
	room, err := r.validateRoomAndPlayer(request.RoomKey, request.PlayerKey, request.PlayerSecret)
	if err != nil {
		return nil, err
	}
	return &rpcRequestContext{
		process: func() (any, error) {
			area := make([]contest.Coordinate, 0, len(request.Area))
			for _, coordinate := range request.Area {
				area = append(
					area, contest.Coordinate{
						Lat: coordinate[0],
						Lng: coordinate[1],
					},
				)
			}
			room.SetOptions(
				contest.RoomOptions{
					Area:              area,
					NumberOfQuestions: request.NumberOfQuestions,
					MaxAnswerTime:     time.Duration(request.MaxAnswerTimeSec) * time.Second,
				}, request.PlayerKey,
			)
			return updateRoomResponse{
				Errors: room.ConfigErrors(),
			}, nil
		},
		release: unlockRoom(room),
	}, nil
}

type joinRequest struct {
	Name    string `json:"name"`
	RoomKey string `json:"roomKey"`
}

type joinResponse struct {
	Name         string `json:"name"`
	PlayerKey    string `json:"playerKey"`
	PlayerSecret string `json:"playerSecret"`
}

func (r *roomContainer) joinRoom(message json.RawMessage) (*rpcRequestContext, error) {
	request := parseMessage[joinRequest](message)
	r.RLock()
	room, ok := r.openRooms[request.RoomKey]
	r.RUnlock()
	if !ok {
		return &rpcRequestContext{release: unlockRoom(room)}, fmt.Errorf(
			"room with key \"%s\" not found", request.RoomKey,
		)
	}
	room.Lock()
	return &rpcRequestContext{
		process: func() (any, error) {
			player := room.Join(request.Name)
			response := joinResponse{
				Name:         player.Name,
				PlayerKey:    player.Key,
				PlayerSecret: player.Secret,
			}
			log.Printf("Player \"%s\" (\"%s\") joined room \"%s\".", player.Secret, player.Name, room.Key())
			return response, nil
		},
		release: unlockRoom(room),
	}, nil
}

type startGameRequest struct {
	PlayerKey    string `json:"playerKey"`
	PlayerSecret string `json:"playerSecret"`
	RoomKey      string `json:"roomKey"`
}

func (r *roomContainer) startGame(message json.RawMessage) (*rpcRequestContext, error) {
	request := parseMessage[startGameRequest](message)
	room, err := r.validateRoomAndPlayer(request.RoomKey, request.PlayerKey, request.PlayerSecret)
	if err != nil {
		return &rpcRequestContext{release: unlockRoom(room)}, err
	}
	if len(room.ConfigErrors()) > 0 {
		return &rpcRequestContext{release: unlockRoom(room)}, fmt.Errorf(
			"the room is not yet configured correctly, "+
				"it still has the following config errors: %v", room.ConfigErrors(),
		)
	}
	return &rpcRequestContext{
		process: func() (any, error) {
			room.Play(request.PlayerKey)
			return map[string]any{}, nil
		},
		release: unlockRoom(room),
	}, nil
}

func (r *roomContainer) validateRoomAndPlayer(roomKey string, playerKey string, secret string) (contest.Room, error) {
	r.RLock()
	room, ok := r.openRooms[roomKey]
	r.RUnlock()
	if !ok {
		return nil, fmt.Errorf("room with key \"%s\" not found", roomKey)
	}
	room.Lock()
	if p, ok := room.FindPlayer(playerKey); !ok || p.Secret != secret {
		return nil, fmt.Errorf("player with key \"%s\" has not joined the room yet or secret is wrong", playerKey)
	}
	return room, nil
}

type guessRequest struct {
	PlayerKey    string     `json:"playerKey"`
	PlayerSecret string     `json:"playerSecret"`
	RoomKey      string     `json:"roomKey"`
	Guess        [2]float64 `json:"guess"`
}

func (r *roomContainer) answerQuestion(message json.RawMessage) (*rpcRequestContext, error) {
	request := parseMessage[guessRequest](message)
	room, err := r.validateRoomAndPlayer(request.RoomKey, request.PlayerKey, request.PlayerSecret)
	if err != nil {
		return &rpcRequestContext{release: unlockRoom(room)}, err
	}
	if !room.HasActiveQuestion() {
		return &rpcRequestContext{release: unlockRoom(room)}, fmt.Errorf("question cannot be answered because there is no active question")
	}
	return &rpcRequestContext{
		process: func() (any, error) {
			result, err := room.AnswerQuestion(
				request.PlayerKey, contest.Coordinate{
					Lat: request.Guess[0],
					Lng: request.Guess[1],
				},
			)
			if err != nil {
				return nil, fmt.Errorf("could not validate answer: %v", err)
			}
			return map[string]int{"points": result}, nil
		},
		release: unlockRoom(room),
	}, nil
}

func (r *roomContainer) advanceGame(message json.RawMessage) (*rpcRequestContext, error) {
	request := parseMessage[startGameRequest](message)
	room, err := r.validateRoomAndPlayer(request.RoomKey, request.PlayerKey, request.PlayerSecret)
	if err != nil {
		return &rpcRequestContext{release: unlockRoom(room)}, err
	}
	if !room.CanBeAdvanced() {
		return &rpcRequestContext{release: unlockRoom(room)}, fmt.Errorf(
			"the room \"%s\" cannot be advanced", request.RoomKey,
		)
	}
	return &rpcRequestContext{
		process: func() (any, error) {
			room.AdvanceToNextQuestion()
			return map[string]any{}, nil
		}, release: unlockRoom(room),
	}, err
}

func parseMessage[K any](message json.RawMessage) K {
	var request K
	_ = json.Unmarshal(message, &request)
	return request
}

func unlockRoom(room contest.Room) func() {
	return func() {
		if room != nil {
			room.Unlock()
		}
	}
}
