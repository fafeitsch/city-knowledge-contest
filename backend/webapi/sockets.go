package webapi

import (
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"github.com/fafeitsch/city-knowledge-contest/backend/types"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strings"
	"time"
)

func (r *roomContainer) upgradeToWebSocket(writer http.ResponseWriter, request *http.Request, options Options) error {
	parts := strings.Split(request.RequestURI, "/")
	r.RLock()
	room, roomExists := r.openRooms[parts[2]]
	r.RUnlock()
	if len(parts) != 4 || !roomExists {
		writer.WriteHeader(http.StatusNotFound)
		return nil
	}
	//secret := request.Header.Get("ckc-player-secret")
	player, _ := room.FindPlayer(parts[3])
	//if !ok || player.Secret != secret {
	//	writer.WriteHeader(http.StatusUnauthorized)
	//	return nil
	//}
	connection, err := websocket.Accept(
		writer, request, &websocket.AcceptOptions{InsecureSkipVerify: options.AllowCors},
	)
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
	players := make([]playerInfo, 0, len(existingPlayers))
	for _, p := range existingPlayers {
		players = append(
			players, playerInfo{
				Name:      p.Name,
				PlayerKey: p.Key,
			},
		)
	}
	_ = wsjson.Write(
		request.Context(), connection, websocketMessage{
			Topic:   "successfullyJoined",
			Payload: initialJoinMessage{Players: players, Options: convertRoomOptions(room.Options(), "")},
		},
	)
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
	Players []playerInfo      `json:"players"`
	Options roomUpdateMessage `json:"options"`
}

type playerInfo struct {
	Name      string `json:"name"`
	PlayerKey string `json:"playerKey"`
}

type roomUpdateMessage struct {
	ListFileName      string     `json:"listFileName"`
	Center            [2]float64 `json:"center"`
	NumberOfQuestions int        `json:"numberOfQuestions"`
	MaxAnswerTimeSec  int        `json:"maxAnswerTimeSec"`
	PlayerKey         string     `json:"playerKey,omitempty"`
	Errors            []string   `json:"errors"`
}

type websocketNotifier struct {
	write func(msg any)
}

func (w *websocketNotifier) NotifyPlayerJoined(name string, key string) {
	payload := map[string]string{"name": name, "playerKey": key}
	w.write(websocketMessage{Topic: "playerJoined", Payload: payload})
}

func (w *websocketNotifier) NotifyRoomUpdated(options contest.RoomOptions, playerKey string) {
	message := convertRoomOptions(options, playerKey)
	w.write(
		websocketMessage{
			Topic:   "roomUpdated",
			Payload: message,
		},
	)
}

func (w *websocketNotifier) NotifyGameStarted(playerKey string, center types.Coordinate) {
	message := map[string]any{"playerKey": playerKey, "center": [2]float64{center.Lat, center.Lng}}
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
		"delta":  result.PointDelta,
		"points": result.Points,
	}
	w.write(websocketMessage{Topic: "questionFinished", Payload: message})
}

func (w *websocketNotifier) NotifyGameEnded(reason string, result map[string]int) {
	message := map[string]any{
		"reason": reason,
		"result": result,
	}
	w.write(websocketMessage{Topic: "gameEnded", Payload: message})
}

func (w *websocketNotifier) NotifyPlayerAnswered(playerKey string) {
	message := map[string]any{"playerKey": playerKey}
	w.write(websocketMessage{Topic: "playerAnswered", Payload: message})
}

func convertRoomOptions(options contest.RoomOptions, playerKey string) roomUpdateMessage {
	listName := ""
	if options.StreetList != nil {
		listName = options.StreetList.FileName
	}
	center := [2]float64{0, 0}
	if options.StreetList != nil {
		center = [2]float64{options.StreetList.Center.Lat, options.StreetList.Center.Lng}
	}
	message := roomUpdateMessage{
		ListFileName:      listName,
		Center:            center,
		MaxAnswerTimeSec:  int(options.MaxAnswerTime / time.Second),
		NumberOfQuestions: options.NumberOfQuestions,
		PlayerKey:         playerKey,
		Errors:            options.Errors(),
	}
	return message
}
