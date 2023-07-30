package webapi

import (
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sync"
	"time"
)

type statisticsContainer struct {
	sync.RWMutex
	subscribers []websocketNotifier
	roomKeys    map[string]string
}

func (s *statisticsContainer) createSocket(
	writer http.ResponseWriter, request *http.Request, options Options,
) error {
	connection, err := websocket.Accept(
		writer, request, &websocket.AcceptOptions{InsecureSkipVerify: options.AllowCors},
	)
	if err != nil {
		return fmt.Errorf("could not upgrade to websockets: %v", err)
	}
	notifier := websocketNotifier{
		write: func(msg any) {
			_ = wsjson.Write(request.Context(), connection, msg)
		},
	}
	log.Printf("connection to statistics subscriber established")
	s.Lock()
	s.subscribers = append(s.subscribers, notifier)
	s.Unlock()
	closeContext := connection.CloseRead(request.Context())
	var pingErr error
	for pingErr == nil {
		pingErr = connection.Ping(closeContext)
		time.Sleep(10 * time.Second)
	}
	log.Printf("Connection to statistics subscriber lost")
	// TODO remove subscriber from list
	_ = connection.Close(websocket.StatusNormalClosure, "")
	return nil
}

func (s *statisticsContainer) pseudomizeRoomKey(roomKey string) string {
	if s.roomKeys[roomKey] == "" {
		s.roomKeys[roomKey] = keygen.RoomKey()
	}
	return s.roomKeys[roomKey]
}

func (s *statisticsContainer) sendRoomCreated(key string) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(map[string]string{"topic": "room created", "roomKey": s.pseudomizeRoomKey(key)})
	}
}

func (s *statisticsContainer) sendPlayerJoined(key string) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(map[string]string{"topic": "player joined", "roomKey": s.pseudomizeRoomKey(key)})
	}
}

func (s *statisticsContainer) sendPlayerLeft(key string) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(map[string]string{"topic": "player left", "roomKey": s.pseudomizeRoomKey(key)})
	}
}

func (s *statisticsContainer) sendGameStarted(key string, options contest.RoomOptions) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(
			map[string]any{
				"topic":               "game started",
				"roomKey":             s.pseudomizeRoomKey(key),
				"questions":           options.NumberOfQuestions,
				"maxAnswerTimeSecond": options.MaxAnswerTime / time.Second,
				"streelist":           options.StreetList.Name,
			},
		)
	}
}

func (s *statisticsContainer) sendQuestionAnswered(key string, points int) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(
			map[string]any{
				"topic":   "question answered",
				"roomKey": s.pseudomizeRoomKey(key),
				"points":  points,
			},
		)
	}
}

func (s *statisticsContainer) sendGameAdvanced(key string) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(map[string]string{"topic": "game advanced", "roomKey": s.pseudomizeRoomKey(key)})
	}
}

func (s *statisticsContainer) sendRoomCleared(key string, finished bool) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(
			map[string]any{
				"topic":    "room cleared",
				"roomKey":  s.pseudomizeRoomKey(key),
				"finished": finished,
			},
		)
	}
}
