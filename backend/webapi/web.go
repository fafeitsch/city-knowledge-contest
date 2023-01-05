package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
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
	AllowCors bool
}

type createRoomRequest struct {
	Name string `json:"name"`
	Area [][2]float64
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

func HandleFunc(options Options) http.HandlerFunc {
	methods := map[string]Method{
		"createRoom": func(message json.RawMessage) (json.RawMessage, error) {
			var request createRoomRequest
			_ = json.Unmarshal(message, &request)
			if len(request.Area) < 3 {
				return nil, fmt.Errorf("the area of the room must consist of at least tree coordinates, "+
					"but consists of %d", len(request.Area))
			}
			area := make([]contest.Coordinate, 0, len(request.Area))
			for _, coordinate := range request.Area {
				area = append(area, contest.Coordinate{
					Lat: coordinate[0],
					Lng: coordinate[1],
				})
			}
			room := contest.NewRoom(area)
			player := room.Join(request.Name)
			openRooms[room.Key] = room
			response := createRoomResponse{RoomKey: room.Key, PlayerKey: player.Key}
			log.Printf("Created room \"%s\" from player \"%s\" (\"%s\").", room.Key, player.Key, player.Name)
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
				Name:      player.Name,
				PlayerKey: player.Key,
			}
			log.Printf("Player \"%s\" (\"%s\") joined room \"%s\".", player.Key, player.Name, room.Key)
			msg, _ := json.Marshal(response)
			return msg, nil
		},
	}
	return func(resp http.ResponseWriter, req *http.Request) {
		if options.AllowCors {
			setCorsHeaders(resp)
		}
		if req.Method == "OPTIONS" {
			return
		}
		if req.Header.Get("Upgrade") == "websocket" {
			err := handleWebsocketRequest(resp, req, options)
			if err != nil {
				fmt.Printf("%v", err)
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

func setCorsHeaders(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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

func handleWebsocketRequest(writer http.ResponseWriter, request *http.Request, options Options) error {
	parts := strings.Split(request.RequestURI, "/")
	room, roomExists := openRooms[parts[2]]
	if len(parts) != 4 || !roomExists {
		writer.WriteHeader(http.StatusNotFound)
	}
	player := room.FindPlayer(parts[3])
	if player == nil {
		writer.WriteHeader(http.StatusNotFound)
	}
	connection, err := websocket.Accept(writer, request, &websocket.AcceptOptions{InsecureSkipVerify: options.AllowCors})
	if err != nil {
		return fmt.Errorf("could not upgrade to websockets: %v", err)
	}
	notifier := &websocketNotifier{
		write: func(msg any) {
			_ = wsjson.Write(request.Context(), connection, msg)
		},
	}
	player.Notifier = notifier
	existingPlayers := room.Players()
	players := make([]string, 0, len(existingPlayers))
	for _, p := range existingPlayers {
		players = append(players, p.Name)
	}
	_ = wsjson.Write(request.Context(), connection, websocketMessage{
		Topic:   "successfullyJoined",
		Payload: initialJoinMessage{Players: players},
	})
	log.Printf("Established websocket connection to player \"%s\" (\"%s\").", player.Key, player.Name)
	alive := true
	for alive {
		time.Sleep(10 * time.Second)
		func() {
			pingContext, cancel := context.WithTimeout(request.Context(), time.Second)
			defer cancel()
			pingErr := connection.Ping(pingContext)
			alive = pingErr == nil
		}()
	}
	log.Printf("Connection to player \"%s\" (\"%s\") lost.", player.Key, player.Name)
	_ = connection.Close(websocket.StatusNormalClosure, "")
	return nil
}

type websocketNotifier struct {
	write func(msg any)
}

func (w *websocketNotifier) NotifyPlayerJoined(name string) {
	payload := map[string]string{"name": name}
	w.write(websocketMessage{Topic: "playerJoined", Payload: payload})
}
