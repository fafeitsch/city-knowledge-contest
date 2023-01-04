package contest

import (
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"sync"
	"time"
)

type Room struct {
	mutex    sync.RWMutex
	Key      string
	Creation time.Time
	players  []Player
}

func (r *Room) Players() []Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.players
}

func NewRoom() *Room {
	return &Room{Key: keygen.Generate(), Creation: time.Now(), players: make([]Player, 0, 0)}
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
	player := Player{Name: name, Key: keygen.Generate()}
	r.players = append(r.players, player)
	return player
}

func (r *Room) FindPlayer(key string) *Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for index, player := range r.players {
		if player.Key == key {
			return &r.players[index]
		}
	}
	return nil
}

type Player struct {
	Notifier
	Key  string
	Name string
}

type Notifier interface {
	NotifyPlayerJoined(name string)
}
