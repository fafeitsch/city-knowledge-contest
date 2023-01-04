package contest

import (
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"time"
)

type Room struct {
	Key      string
	Creation time.Time
	Players  []Player
}

func NewRoom() *Room {
	return &Room{Key: keygen.Generate(), Creation: time.Now(), Players: make([]Player, 0, 0)}
}

type Player struct {
	Key  string
	Name string
}

func (r *Room) Join(name string) Player {
	player := Player{Name: name, Key: keygen.Generate()}
	r.Players = append(r.Players, player)
	return player
}
