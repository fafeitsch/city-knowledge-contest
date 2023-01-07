package webapi

import (
	"encoding/json"
	"reflect"
)

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Id      *string         `json:"id"`
}

type Response struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
	Id      *string         `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Handler interface {
	Methods() map[string]rpcMethod
}

type rpcMethod struct {
	description    string
	input          reflect.Type
	output         reflect.Type
	method         Method
	persistChanged bool
}

type Method func(message json.RawMessage) (json.RawMessage, error)

type Options struct {
	AllowCors bool
}

type createRoomRequest struct {
	Name string `json:"name"`
}

type createRoomResponse struct {
	RoomKey   string `json:"roomKey"`
	PlayerKey string `json:"playerKey"`
}

type joinRequest struct {
	Name    string `json:"name"`
	RoomKey string `json:"roomKey"`
}

type joinResponse struct {
	Name      string `json:"name"`
	PlayerKey string `json:"playerKey"`
}

type websocketMessage struct {
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}

type initialJoinMessage struct {
	Players []string          `json:"players"`
	Options roomUpdateMessage `json:"options"`
}

type roomUpdateRequest struct {
	Area              [][2]float64 `json:"area"`
	NumberOfQuestions int          `json:"numberOfQuestions"`
	RoomKey           string       `json:"roomKey"`
	PlayerKey         string       `json:"playerKey"`
}

type roomUpdateMessage struct {
	Area              [][2]float64 `json:"area"`
	NumberOfQuestions int          `json:"numberOfQuestions"`
	PlayerName        string       `json:"playerName,omitempty"`
	Errors            []string     `json:"errors"`
}

type startGameRequest struct {
	PlayerKey string `json:"playerKey"`
	RoomKey   string `json:"roomKey"`
}

type guessRequest struct {
	PlayerKey string     `json:"playerKey"`
	RoomKey   string     `json:"roomKey"`
	Guess     [2]float64 `json:"guess"`
}
