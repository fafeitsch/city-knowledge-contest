package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"runtime/debug"
	"strings"
	"time"
)

var openRooms = make(map[string]*contest.Room)

func parseMessage[K any](message json.RawMessage) K {
	var request K
	_ = json.Unmarshal(message, &request)
	return request
}

func HandleFunc(options Options) http.HandlerFunc {
	methods := map[string]Method{
		"createRoom": func(message json.RawMessage) (json.RawMessage, error) {
			request := parseMessage[createRoomRequest](message)
			room := contest.NewRoom()
			player := room.Join(request.Name)
			openRooms[room.Key] = room
			response := createRoomResponse{
				RoomKey:   room.Key,
				PlayerKey: player.Key,
			}
			log.Printf("Created room \"%s\" from player \"%s\" (\"%s\").", room.Key, player.Key, player.Name)
			msg, _ := json.Marshal(response)
			return msg, nil
		},
		"updateRoom": func(message json.RawMessage) (json.RawMessage, error) {
			request := parseMessage[roomUpdateRequest](message)
			room, ok := openRooms[request.RoomKey]
			if !ok {
				return nil, fmt.Errorf("room with key \"%s\" not found", request.RoomKey)
			}
			if room.FindPlayer(request.PlayerKey) == nil {
				return nil, fmt.Errorf(
					"player with key \"%s\" has not joined the room yet",
					request.PlayerKey)
			}
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
			return []byte("{}"), nil
		},
		"joinRoom": func(message json.RawMessage) (json.RawMessage, error) {
			request := parseMessage[joinRequest](message)
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
		"startGame": func(message json.RawMessage) (json.RawMessage, error) {
			request := parseMessage[startGameRequest](message)
			room, ok := openRooms[request.RoomKey]
			if !ok {
				return nil, fmt.Errorf("room with key \"%s\" not found", request.RoomKey)
			}
			if room.FindPlayer(request.PlayerKey) == nil {
				return nil, fmt.Errorf(
					"player with key \"%s\" has not joined the room yet",
					request.PlayerKey)
			}
			go room.Play(request.PlayerKey)
			return []byte("{}"), nil
		},
		"answerQuestion": func(message json.RawMessage) (json.RawMessage, error) {
			request := parseMessage[guessRequest](message)
			room, ok := openRooms[request.RoomKey]
			if !ok {
				return nil, fmt.Errorf("room with key \"%s\" not found", request.RoomKey)
			}
			if room.FindPlayer(request.PlayerKey) == nil {
				return nil, fmt.Errorf(
					"player with key \"%s\" has not joined the room yet",
					request.PlayerKey)
			}
			result, err := room.AnswerQuestion(request.PlayerKey, contest.Coordinate{
				Lat: request.Guess[0],
				Lng: request.Guess[1],
			})
			if err != nil {
				return nil, fmt.Errorf("could not validate answer: %v", err)
			}
			return []byte(fmt.Sprintf("{\"correct\": %v}", result)), nil
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
		Payload: initialJoinMessage{Players: players, Options: convertRoomOptions(room.Options(), "")},
	})
	log.Printf("Established websocket connection to player \"%s\" (\"%s\").", player.Key, player.Name)
	closeContext := connection.CloseRead(request.Context())
	var pingErr error
	for pingErr == nil {
		pingErr = connection.Ping(closeContext)
		time.Sleep(10 * time.Second)
	}
	log.Printf("Connection to player \"%s\" (\"%s\") lost: %v", player.Key, player.Name, pingErr)
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

func (w *websocketNotifier) NotifyRoomUpdated(options contest.RoomOptions, playerName string) {
	message := convertRoomOptions(options, playerName)
	w.write(websocketMessage{
		Topic:   "roomUpdated",
		Payload: message,
	})
}

func (w *websocketNotifier) NotifyGameStarted(player string, center contest.Coordinate) {
	message := map[string]any{"playerName": player, "center": [2]float64{center.Lat, center.Lng}}
	w.write(websocketMessage{Topic: "gameStarted", Payload: message})
}

func (w *websocketNotifier) NotifyQuestionCountdown(followUps int) {
	message := map[string]any{"followUps": followUps}
	w.write(websocketMessage{Topic: "questionCountdown", Payload: message})
}

func (w *websocketNotifier) NotifyQuestion(question string) {
	message := map[string]string{"find": question}
	w.write(websocketMessage{Topic: "question", Payload: message})
}

func (w *websocketNotifier) NotifyAnswerTimeCountdown(followUps int) {
	message := map[string]any{"followUps": followUps}
	w.write(websocketMessage{Topic: "answerCountdown", Payload: message})
}

func (w *websocketNotifier) NotifyQuestionResults(result contest.QuestionResult) {
	message := map[string]any{
		"question": result.Question,
		"solution": [2]float64{
			result.Solution.Lat,
			result.Solution.Lng,
		},
	}
	w.write(websocketMessage{Topic: "questionFinished", Payload: message})
}

func convertRoomOptions(options contest.RoomOptions, playerName string) roomUpdateMessage {
	area := make([][2]float64, 0, len(options.Area))
	for _, coordinate := range options.Area {
		area = append(area, [2]float64{coordinate.Lat, coordinate.Lng})
	}
	message := roomUpdateMessage{
		Area:              area,
		NumberOfQuestions: options.NumberOfQuestions,
		PlayerName:        playerName,
		Errors:            options.Errors(),
	}
	return message
}
