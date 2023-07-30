package webapi

import (
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strings"
	"sync"
	"time"
)

type statisticsContainer struct {
	sync.RWMutex
	secretUrlKey string
	subscribers  []websocketNotifier
}

func (s *statisticsContainer) createSocket(
	writer http.ResponseWriter, request *http.Request, options Options,
) error {
	parts := strings.Split(request.RequestURI, "/")
	if parts[2] != s.secretUrlKey {
		writer.WriteHeader(http.StatusNotFound)
		return nil
	}
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

func (s *statisticsContainer) sendRoomCreated(key string) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(map[string]string{"topic": "room created", "roomKey": key})
	}
}

func (s *statisticsContainer) sendRoomJoined(key string) {
	s.RLock()
	defer s.RUnlock()
	for _, notifier := range s.subscribers {
		notifier.write(map[string]string{"topic": "room joined", "roomKey": key})
	}
}
