package contest

import (
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"sync"
	"time"
)

type Coordinate struct {
	Lat float64
	Lng float64
}

type Room struct {
	mutex    sync.RWMutex
	Key      string
	Creation time.Time
	players  []Player
	options  RoomOptions
}

type RoomOptions struct {
	Area              []Coordinate
	NumberOfQuestions int
}

func (r *RoomOptions) Errors() []string {
	errors := make([]string, 0, 0)
	if len(r.Area) < 3 {
		errors = append(errors, "toFewCoordinates")
	}
	if r.NumberOfQuestions < 1 {
		errors = append(errors, "numberOfQuestionsToSmall")
	}
	if r.NumberOfQuestions > 100 {
		errors = append(errors, "numberOfQuestionsToBig")
	}
	return errors
}

func (r *Room) Players() []Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.players
}

func NewRoom() *Room {
	return &Room{
		Key:      keygen.RoomKey(),
		Creation: time.Now(),
		players:  make([]Player, 0, 0),
		options: RoomOptions{
			Area:              []Coordinate{},
			NumberOfQuestions: 10,
		},
	}
}

func (r *Room) Options() RoomOptions {
	return r.options
}

func (r *Room) SetOptions(options RoomOptions, playerKey string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	player := r.findPlayer(playerKey)
	if player == nil {
		panic(fmt.Sprintf("player with key \"%s\" does not exist in room \"%s\"", playerKey, r.Key))
	}
	r.options = options
	for _, p := range r.players {
		if p.Notifier == nil {
			continue
		}
		p.NotifyRoomUpdated(options, player.Name)
	}
}

func (r *Room) Join(name string) Player {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, player := range r.players {
		if player.Notifier == nil {
			continue
		}
		player.NotifyPlayerJoined(name)
	}
	player := Player{Name: name, Key: keygen.PlayerKey()}
	r.players = append(r.players, player)
	return player
}

func (r *Room) findPlayer(key string) *Player {
	for index, player := range r.players {
		if player.Key == key {
			return &r.players[index]
		}
	}
	return nil
}

func (r *Room) FindPlayer(key string) *Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.findPlayer(key)
}

type Player struct {
	Notifier
	Key  string
	Name string
}

type Notifier interface {
	NotifyPlayerJoined(string)
	NotifyRoomUpdated(RoomOptions, string)
}
