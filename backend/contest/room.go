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
	mutex           sync.Mutex
	Key             string
	Creation        time.Time
	players         map[string]*Player
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
		errors = append(errors, "tooFewCoordinates")
	}
	if r.NumberOfQuestions < 1 {
		errors = append(errors, "numberOfQuestionsToSmall")
	}
	if r.NumberOfQuestions > 100 {
		errors = append(errors, "numberOfQuestionsToBig")
	}
	return errors
}

func (r *Room) Lock() {
	r.mutex.Lock()
}

func (r *Room) Unlock() {
	r.mutex.Unlock()
}

func (r *Room) Players() []Player {
	result := make([]Player, 0, len(r.players))
	for key := range r.players {
		result = append(result, *r.players[key])
	}
	return result
}

func NewRoom() *Room {
	return &Room{
		Key:      keygen.RoomKey(),
		Creation: time.Now(),
		random:   rand.New(rand.NewSource(time.Now().Unix())),
		players:  make(map[string]*Player),
		options: RoomOptions{
			Area:              []Coordinate{},
			NumberOfQuestions: 10,
		},
	}
}

func (r *Room) Options() RoomOptions {
	return r.options
}

func (r *Room) ConfigErrors() []string {
	return r.options.Errors()
}

func (r *Room) SetOptions(options RoomOptions, playerKey string) {
	player, ok := r.players[playerKey]
	if !ok {
		panic(fmt.Sprintf("player with key \"%s\" does not exist in room \"%s\"", playerKey, r.Key))
	}
	r.options = options
	r.notifyPlayers(func(p Player) {
		p.NotifyRoomUpdated(options, player.Name)
	})
}

func (r *Room) Join(name string) Player {
	player := Player{Name: name, Secret: keygen.PlayerKey(), Key: keygen.PlayerKey()}
	r.notifyPlayers(func(p Player) {
		p.NotifyPlayerJoined(player.Name, player.Key)
	})
	r.players[player.Key] = &player
	return player
}

func (r *Room) FindPlayer(key string) (*Player, bool) {
	result, ok := r.players[key]
	return result, ok
}

func (r *Room) notifyPlayers(consumer func(Player)) {
	for _, player := range r.players {
		if player.Notifier == nil {
			continue
		}
		consumer(*player)
	}
}

func (r *Room) Play(playerKey string) {
	if len(r.options.Errors()) > 0 {
		panic(fmt.Sprintf("can't start the game because there are still errors in its config"))
	}
	startPlayer, ok := r.players[playerKey]
	if !ok {
		panic(fmt.Sprintf("player with key \"%s\" does not exist in this room", playerKey))
	}
	triangles := polygon(r.Options().Area).computeTriangulation()
	center := triangles.randomPoint(r.random)
	r.notifyPlayers(func(player Player) {
		player.NotifyGameStarted(startPlayer.Name, center)
	})
	numberOfQuestions := r.options.NumberOfQuestions
	go func() {
		for round := 0; round < numberOfQuestions; round++ {
			r.playQuestion(round, triangles)
		}
	}()
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
	r.Lock()
	r.currentQuestion = &Question{
		StreetName:         question,
		Solution:           solution,
		points:             make(map[string]int),
		allPlayersAnswered: make(chan bool),
	}
	r.Unlock()
	r.notifyPlayers(func(player Player) {
		player.NotifyQuestion(question)
	})
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
	r.Lock()
	r.currentQuestion = nil
	r.Unlock()

	return nil
}

func (r *Room) AnswerQuestion(playerKey string, guess Coordinate) (bool, error) {
	_, ok := r.players[playerKey]
	if !ok {
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

func (r *Room) HasActiveQuestion() bool {
	return r.currentQuestion != nil
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
	Secret string
	Key    string
	Name   string
}

type QuestionResult struct {
	Question   string
	Solution   Coordinate
	PointDelta map[string]int
	Points     map[string]int
}

type Notifier interface {
	NotifyPlayerJoined(string, string)
	NotifyRoomUpdated(RoomOptions, string)
	NotifyGameStarted(string, Coordinate)
	NotifyQuestionCountdown(int)
	NotifyQuestion(string)
	NotifyAnswerTimeCountdown(int)
	NotifyQuestionResults(result QuestionResult)
}
