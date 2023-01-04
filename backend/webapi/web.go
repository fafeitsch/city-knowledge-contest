package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"math/rand"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"reflect"
	"runtime/debug"
	"strings"
)

var openRooms = make(map[string]*contest.Room)

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
	Random *rand.Rand
}

type createRoomResponse struct {
	Key string `json:"key"`
}

type joinRequest struct {
	Name    string `json:"name"`
	RoomKey string `json:"roomKey"`
}

type joinResponse struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func HandleFunc(options Options) http.HandlerFunc {
	methods := map[string]Method{
		"createRoom": func(message json.RawMessage) (json.RawMessage, error) {
			room := contest.NewRoom()
			openRooms[room.Key] = room
			response := createRoomResponse{Key: room.Key}
			msg, _ := json.Marshal(response)
			return msg, nil
		},
		"joinRoom": func(message json.RawMessage) (json.RawMessage, error) {
			var request joinRequest
			_ = json.Unmarshal(message, &request)
			room, ok := openRooms[request.RoomKey]
			if !ok {
				return nil, fmt.Errorf("room with key \"%s\" not found", request.RoomKey)
			}
			player := room.Join(request.Name)
			response := joinResponse{
				Name: player.Name,
				Key:  player.Key,
			}
			msg, _ := json.Marshal(response)
			return msg, nil
		},
	}
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if req.Method == "OPTIONS" {
			return
		}
		if req.Header.Get("Upgrade") == "websocket" {
			err := handleWebsocketRequest(resp, req)
			if err != nil {
				fmt.Printf("%v")
				resp.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		if req.Method != "POST" {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(req.RequestURI, "/")
		if len(parts) != 2 || parts[1] != "rpc" {
			resp.WriteHeader(http.StatusNotFound)
			return
		}
		var request Request
		defer handlePanic(resp, request)
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			writeError(resp, -32700, nil, "could not parse JSON-RPC request: %v", err)
			return
		}
		method, ok := methods[request.Method]
		if !ok {
			writeError(resp, -32601, request.Id, "the requested method \"%s\" was not found", request.Method)
			return
		}
		result, err := method(request.Params)
		if err != nil {
			writeError(resp, -32603, request.Id, "the method \"%s\" could not be executed properly: %v", request.Method, err)
			return
		}
		response := Response{Id: request.Id, Jsonrpc: "2.0", Result: result}
		_ = json.NewEncoder(resp).Encode(response)
	}
}

func handlePanic(resp http.ResponseWriter, request Request) {
	if r := recover(); r != nil {
		writeError(resp, -32603, request.Id, "a fatal error occurred during the method execution.")
		fmt.Printf("panic during RPC call: %v\n%v", r, string(debug.Stack()))
	}
}

func writeError(resp http.ResponseWriter, code int, id *string, format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)
	response := Response{
		Id:      id,
		Error:   &Error{Code: code, Message: msg},
		Jsonrpc: "2.0",
	}
	resp.WriteHeader(400)
	_ = json.NewEncoder(resp).Encode(response)
}

type websocketMessage struct {
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}

type initialJoinMessage struct {
	Players []string `json:"players"`
}

func handleWebsocketRequest(writer http.ResponseWriter, request *http.Request) error {
	parts := strings.Split(request.RequestURI, "/")
	room, roomExists := openRooms[parts[2]]
	if len(parts) != 4 || !roomExists {
		writer.WriteHeader(http.StatusNotFound)
	}
	var player *contest.Player = nil
	players := make([]string, 0, len(room.Players))
	for _, p := range room.Players {
		if p.Key == parts[3] {
			player = &p
		}
		players = append(players, p.Name)
	}
	if player == nil {
		writer.WriteHeader(http.StatusNotFound)
	}
	connection, err := websocket.Accept(writer, request, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return fmt.Errorf("could not upgrade to websockets: %v", err)
	}
	wsjson.Write(request.Context(), connection, websocketMessage{Topic: "successfullyJoined", Payload: initialJoinMessage{Players: players}})
	connection.Close(websocket.StatusNormalClosure, "Nothing more implemented yet.")
	return nil
}
