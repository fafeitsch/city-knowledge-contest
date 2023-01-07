package contest

import (
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"math/rand"
	"sync"
	"time"
)

type Coordinate struct {
	Lat float64
	Lng float64
}

type Room struct {
	mutex           sync.RWMutex
	Key             string
	Creation        time.Time
	players         []Player
	random          *rand.Rand
	options         RoomOptions
	currentQuestion *Question
}

type RoomOptions struct {
	Area              []Coordinate
	NumberOfQuestions int
}

type Question struct {
	StreetName         string
	Solution           Coordinate
	points             map[string]int
	allPlayersAnswered chan bool
}

func (q *Question) waitForPlayers(countdown func(int)) {
	questionTimeout := 120 * time.Second
	t2 := time.NewTimer(questionTimeout - (1 * time.Second))
	t1 := time.NewTimer(questionTimeout - (2 * time.Second))
	finish := time.NewTimer(questionTimeout)
	finished := false
	for !finished {
		select {
		case <-q.allPlayersAnswered:
			t1.Stop()
			t2.Stop()
			finish.Stop()
			finished = true
		case <-t1.C:
			countdown(1)
		case <-t2.C:
			countdown(0)
			t1.Stop()
		case <-finish.C:
			t1.Stop()
			t2.Stop()
			finished = true
		}
	}
}

func (r *RoomOptions) Errors() []string {
	errors := make([]string, 0, 0)
	if len(r.Area) < 3 {
		errors = append(errors, "toFewCoordinates")
	}
	if r.NumberOfQuestions < 1 {
		errors = append(errors, "numberOfQuestionsToSmall")
	}
	if r.NumberOfQuestions > 100 {
		errors = append(errors, "numberOfQuestionsToBig")
	}
	return errors
}

func (r *Room) Players() []Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.players
}

func NewRoom() *Room {
	return &Room{
		Key:      keygen.RoomKey(),
		Creation: time.Now(),
		random:   rand.New(rand.NewSource(time.Now().Unix())),
		players:  make([]Player, 0, 0),
		options: RoomOptions{
			Area:              []Coordinate{},
			NumberOfQuestions: 10,
		},
	}
}

func (r *Room) Options() RoomOptions {
	return r.options
}

func (r *Room) SetOptions(options RoomOptions, playerKey string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	player := r.findPlayer(playerKey)
	if player == nil {
		panic(fmt.Sprintf("player with key \"%s\" does not exist in room \"%s\"", playerKey, r.Key))
	}
	r.options = options
	r.notifyPlayers(func(p Player) {
		p.NotifyRoomUpdated(options, player.Name)
	})
}

func (r *Room) Join(name string) Player {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.notifyPlayers(func(player Player) {
		player.NotifyPlayerJoined(name)
	})
	player := Player{Name: name, Key: keygen.PlayerKey()}
	r.players = append(r.players, player)
	return player
}

func (r *Room) findPlayer(key string) *Player {
	for index, player := range r.players {
		if player.Key == key {
			return &r.players[index]
		}
	}
	return nil
}

func (r *Room) FindPlayer(key string) *Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.findPlayer(key)
}

func (r *Room) notifyPlayers(consumer func(Player)) {
	for _, player := range r.players {
		if player.Notifier == nil {
			continue
		}
		consumer(player)
	}
}

func (r *Room) Play(playerKey string) {
	if len(r.options.Errors()) > 0 {
		panic(fmt.Sprintf("can't start the game because there are still errors in its config"))
	}
	startPlayer := r.findPlayer(playerKey)
	if startPlayer == nil {
		panic(fmt.Sprintf("player with key \"%s\" does not exist in this room", playerKey))
	}
	triangles := polygon(r.Options().Area).computeTriangulation()
	center := triangles.randomPoint(r.random)
	r.notifyPlayers(func(player Player) {
		player.NotifyGameStarted(startPlayer.Name, center)
	})
	numberOfQuestions := r.options.NumberOfQuestions
	for round := 0; round < numberOfQuestions; round++ {
		r.playQuestion(round, triangles)
	}
}

func (r *Room) playQuestion(round int, triangles triangulation) error {
	question := ""
	solution := Coordinate{}
	tries := 0
	for tries < 10 && question == "" {
		randomPoint := triangles.randomPoint(r.random)
		var err error
		question, solution, err = queryStreet(randomPoint, r.random)
		if err != nil {
			return fmt.Errorf("could not query address of %v: %v", randomPoint, err)
		}
		tries = tries + 1
	}
	if question == "" {
		return fmt.Errorf("could not find random streetname after %d attempts, giving up", tries)
	}
	r.sendCountdowns(3, func(player Player) func(int) {
		return player.NotifyQuestionCountdown
	})
	r.notifyPlayers(func(player Player) {
		player.NotifyQuestion(question)
	})
	r.currentQuestion = &Question{
		StreetName:         question,
		Solution:           solution,
		points:             make(map[string]int),
		allPlayersAnswered: make(chan bool),
	}
	r.currentQuestion.waitForPlayers(func(followUps int) {
		r.notifyPlayers(func(player Player) {
			player.NotifyAnswerTimeCountdown(followUps)
		})
	})
	result := QuestionResult{
		Question:   question,
		Solution:   solution,
		PointDelta: map[string]int{},
		Points:     map[string]int{},
	}
	r.notifyPlayers(func(player Player) {
		player.NotifyQuestionResults(result)
	})

	return nil
}

func (r *Room) AnswerQuestion(playerKey string, guess Coordinate) (bool, error) {
	player := r.findPlayer(playerKey)
	if player == nil {
		panic(fmt.Sprintf("player with key \"%s\" not found in this room", playerKey))
	}
	question := r.currentQuestion
	result, err := verifyAnswer(guess, question.StreetName)
	if result {
		question.points[playerKey] = 100
	} else {
		question.points[playerKey] = 0
	}
	if len(question.points) == len(r.players) {
		question.allPlayersAnswered <- true
	}
	return result, err
}

func (r *Room) sendCountdowns(amount int, consumer func(Player) func(int)) {
	for i := 0; i < amount; i++ {
		r.notifyPlayers(func(player Player) {
			consumer(player)(amount - i - 1)
		})
		time.Sleep(time.Second)
	}
}

type Player struct {
	Notifier
	Key  string
	Name string
}

type QuestionResult struct {
	Question   string
	Solution   Coordinate
	PointDelta map[string]int
	Points     map[string]int
}

type Notifier interface {
	NotifyPlayerJoined(string)
	NotifyRoomUpdated(RoomOptions, string)
	NotifyGameStarted(string, Coordinate)
	NotifyQuestionCountdown(int)
	NotifyQuestion(string)
	NotifyAnswerTimeCountdown(int)
	NotifyQuestionResults(result QuestionResult)
}
