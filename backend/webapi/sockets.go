package webapi

import (
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strings"
	"time"
)

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

type websocketMessage struct {
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}

type initialJoinMessage struct {
	Players []string          `json:"players"`
	Options roomUpdateMessage `json:"options"`
}

type roomUpdateMessage struct {
	Area              [][2]float64 `json:"area"`
	NumberOfQuestions int          `json:"numberOfQuestions"`
	PlayerName        string       `json:"playerName,omitempty"`
	Errors            []string     `json:"errors"`
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
