package contest

import (
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Coordinate struct {
	Lat float64
	Lng float64
}

type Room interface {
	Lock()
	Unlock()
	Players() []Player
	Options() RoomOptions
	ConfigErrors() []string
	SetOptions(RoomOptions, string)
	Join(string) Player
	FindPlayer(string) (*Player, bool)
	Play(string)
	AnswerQuestion(string, Coordinate) (int, error)
	HasActiveQuestion() bool
	CanBeAdvanced() bool
	AdvanceToNextQuestion()
	Key() string
}

type roomImpl struct {
	mutex           sync.Mutex
	key             string
	Creation        time.Time
	players         map[string]*Player
	points          map[string]int
	random          *rand.Rand
	options         RoomOptions
	currentQuestion *Question
	advanceGame     chan bool
}

func (r *roomImpl) Key() string {
	return r.key
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
	begin              time.Time
	duration           time.Duration
}

func (q *Question) waitForPlayers(questionTimeout time.Duration, countdown func(int)) {
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

func (r *roomImpl) Lock() {
	r.mutex.Lock()
}

func (r *roomImpl) Unlock() {
	r.mutex.Unlock()
}

func (r *roomImpl) Players() []Player {
	result := make([]Player, 0, len(r.players))
	for key := range r.players {
		result = append(result, *r.players[key])
	}
	return result
}

func NewRoom() Room {
	return &roomImpl{
		key:      keygen.RoomKey(),
		Creation: time.Now(),
		random:   rand.New(rand.NewSource(time.Now().Unix())),
		players:  make(map[string]*Player),
		options: RoomOptions{
			Area:              []Coordinate{},
			NumberOfQuestions: 10,
		},
	}
}

func (r *roomImpl) Options() RoomOptions {
	return r.options
}

func (r *roomImpl) ConfigErrors() []string {
	return r.options.Errors()
}

func (r *roomImpl) SetOptions(options RoomOptions, playerKey string) {
	player, ok := r.players[playerKey]
	if !ok {
		panic(fmt.Sprintf("player with key \"%s\" does not exist in room \"%s\"", playerKey, r.Key))
	}
	r.options = options
	r.notifyPlayers(func(p Player) {
		p.NotifyRoomUpdated(options, player.Name)
	})
}

func (r *roomImpl) Join(name string) Player {
	player := Player{Name: name, Secret: keygen.PlayerKey(), Key: keygen.PlayerKey()}
	r.notifyPlayers(func(p Player) {
		p.NotifyPlayerJoined(player.Name, player.Key)
	})
	r.players[player.Key] = &player
	return player
}

func (r *roomImpl) FindPlayer(key string) (*Player, bool) {
	result, ok := r.players[key]
	return result, ok
}

func (r *roomImpl) notifyPlayers(consumer func(Player)) {
	for _, player := range r.players {
		if player.Notifier == nil {
			continue
		}
		consumer(*player)
	}
}

func (r *roomImpl) Play(playerKey string) {
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
	r.points = make(map[string]int)
	r.advanceGame = make(chan bool)
	go func() {
		for round := 0; round < numberOfQuestions; round++ {
			err := r.playQuestion(round, triangles)
			if err != nil {
				log.Printf("game in room \"%s\" ended due to an error: %v", r.Key, err)
				r.notifyPlayers(func(player Player) {
					player.NotifyGameEnded("error", r.points)
				})
			}
			if round != numberOfQuestions-1 {
				<-r.advanceGame
			}
		}
		r.notifyPlayers(func(player Player) {
			player.NotifyGameEnded("finished", r.points)
		})
		r.advanceGame = nil
		r.points = nil
	}()
}

func (r *roomImpl) playQuestion(round int, triangles triangulation) error {
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
		begin:              time.Now(),
		duration:           120 * time.Second,
	}
	r.Unlock()
	r.notifyPlayers(func(player Player) {
		player.NotifyQuestion(question)
	})
	r.currentQuestion.waitForPlayers(r.currentQuestion.duration, func(followUps int) {
		r.notifyPlayers(func(player Player) {
			player.NotifyAnswerTimeCountdown(followUps)
		})
	})
	for key, value := range r.currentQuestion.points {
		r.points[key] = r.points[key] + value
	}
	result := QuestionResult{
		Question:   question,
		Solution:   solution,
		PointDelta: r.currentQuestion.points,
		Points:     r.points,
	}
	r.notifyPlayers(func(player Player) {
		player.NotifyQuestionResults(result)
	})
	r.Lock()
	r.currentQuestion = nil
	r.Unlock()

	return nil
}

func (r *roomImpl) AnswerQuestion(playerKey string, guess Coordinate) (int, error) {
	_, ok := r.players[playerKey]
	if !ok {
		panic(fmt.Sprintf("player with key \"%s\" not found in this room", playerKey))
	}
	question := r.currentQuestion
	result, err := verifyAnswer(guess, question.StreetName)
	if result {
		difference := time.Now().Sub(question.begin)
		percent := 1.0 * float64(difference.Milliseconds()) / float64(question.duration.Milliseconds())
		question.points[playerKey] = int(math.Max(10, 100-(100*percent)))
	} else {
		question.points[playerKey] = 0
	}
	if len(question.points) == len(r.players) {
		question.allPlayersAnswered <- true
	}
	return question.points[playerKey], err
}

func (r *roomImpl) HasActiveQuestion() bool {
	return r.currentQuestion != nil
}

func (r *roomImpl) CanBeAdvanced() bool {
	return r.advanceGame != nil
}

func (r *roomImpl) AdvanceToNextQuestion() {
	r.advanceGame <- true
}

func (r *roomImpl) sendCountdowns(amount int, consumer func(Player) func(int)) {
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
	NotifyGameEnded(reason string, result map[string]int)
}
