package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"log"
)

type rpcHandler func(message json.RawMessage) (func() (any, error), error)

type createRoomRequest struct {
	Name string `json:"name"`
}

type createRoomResponse struct {
	RoomKey   string `json:"roomKey"`
	PlayerKey string `json:"playerKey"`
}

func createRoom(message json.RawMessage) (func() (any, error), error) {
	request := parseMessage[createRoomRequest](message)
	if len(request.Name) == 0 {
		return nil, fmt.Errorf("a player name must not be empty")
	}
	return func() (any, error) {
		room := contest.NewRoom()
		player := room.Join(request.Name)
		openRooms[room.Key] = room
		response := createRoomResponse{
			RoomKey:   room.Key,
			PlayerKey: player.Key,
		}
		log.Printf("Created room \"%s\" from player \"%s\" (\"%s\").", room.Key, player.Key, player.Name)
		return response, nil
	}, nil
}

type roomUpdateRequest struct {
	Area              [][2]float64 `json:"area"`
	NumberOfQuestions int          `json:"numberOfQuestions"`
	RoomKey           string       `json:"roomKey"`
	PlayerKey         string       `json:"playerKey"`
}

func updateRoom(message json.RawMessage) (func() (any, error), error) {
	request := parseMessage[roomUpdateRequest](message)
	room, err := validateRoomAndPlayer(request.RoomKey, request.PlayerKey)
	if err != nil {
		return nil, err
	}
	return func() (any, error) {
		area := make([]contest.Coordinate, 0, len(request.Area))
		for _, coordinate := range request.Area {
			area = append(area, contest.Coordinate{
				Lat: coordinate[0],
				Lng: coordinate[1],
			})
		}
		room.SetOptions(contest.RoomOptions{
			Area:              area,
			NumberOfQuestions: request.NumberOfQuestions,
		}, request.PlayerKey)
		return map[string]any{}, nil
	}, nil
}

type joinRequest struct {
	Name    string `json:"name"`
	RoomKey string `json:"roomKey"`
}

type joinResponse struct {
	Name      string `json:"name"`
	PlayerKey string `json:"playerKey"`
}

func joinRoom(message json.RawMessage) (func() (any, error), error) {
	request := parseMessage[joinRequest](message)
	room, ok := openRooms[request.RoomKey]
	if !ok {
		return nil, fmt.Errorf("room with key \"%s\" not found", request.RoomKey)
	}
	return func() (any, error) {
		player := room.Join(request.Name)
		response := joinResponse{
			Name:      player.Name,
			PlayerKey: player.Key,
		}
		log.Printf("Player \"%s\" (\"%s\") joined room \"%s\".", player.Key, player.Name, room.Key)
		return response, nil
	}, nil
}

type startGameRequest struct {
	PlayerKey string `json:"playerKey"`
	RoomKey   string `json:"roomKey"`
}

func startGame(message json.RawMessage) (func() (any, error), error) {
	request := parseMessage[startGameRequest](message)
	room, err := validateRoomAndPlayer(request.RoomKey, request.PlayerKey)
	if err != nil {
		return nil, err
	}
	return func() (any, error) {
		go room.Play(request.PlayerKey)
		return map[string]any{}, nil
	}, nil
}

func validateRoomAndPlayer(roomKey string, playerKey string) (*contest.Room, error) {
	room, ok := openRooms[roomKey]
	if !ok {
		return nil, fmt.Errorf("room with key \"%s\" not found", roomKey)
	}
	if room.FindPlayer(playerKey) == nil {
		return nil, fmt.Errorf("player with key \"%s\" has not joined the room yet", playerKey)
	}
	return room, nil
}

type guessRequest struct {
	PlayerKey string     `json:"playerKey"`
	RoomKey   string     `json:"roomKey"`
	Guess     [2]float64 `json:"guess"`
}

func answerQuestion(message json.RawMessage) (func() (any, error), error) {
	request := parseMessage[guessRequest](message)
	room, err := validateRoomAndPlayer(request.RoomKey, request.PlayerKey)
	if err != nil {
		return nil, err
	}
	return func() (any, error) {
		result, err := room.AnswerQuestion(request.PlayerKey, contest.Coordinate{
			Lat: request.Guess[0],
			Lng: request.Guess[1],
		})
		if err != nil {
			return nil, fmt.Errorf("could not validate answer: %v", err)
		}
		return map[string]bool{"correct": result}, nil
	}, nil
}
