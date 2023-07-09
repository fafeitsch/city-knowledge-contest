package contest

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/fafeitsch/city-knowledge-contest/backend/geodata"
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"github.com/fafeitsch/city-knowledge-contest/backend/types"
)

type Room interface {
	Lock()
	Unlock()
	Players() []Player
	Options() RoomOptions
	ConfigErrors() []string
	SetOptions(RoomOptions, string)
	Join(string) Player
	Leave(string) Player
	FindPlayer(string) (*Player, bool)
	Play(string)
	AnswerQuestion(string, types.Coordinate) (int, error)
	HasActiveQuestion(string) bool
	Question() (string, int)
	CanBeAdvanced() bool
	AdvanceToNextQuestion()
	Key() string
	Finished() bool
	Creation() time.Time
	Started() bool
}

type roomImpl struct {
	mutex           sync.Mutex
	key             string
	creation        time.Time
	players         map[string]*Player
	points          map[string]int
	random          *rand.Rand
	options         RoomOptions
	currentQuestion *Question
	advanceGame     chan bool
	finished        bool
	started         bool
}

func (r *roomImpl) Key() string {
	return r.key
}

type RoomOptions struct {
	StreetList        *geodata.StreetList
	NumberOfQuestions int
	MaxAnswerTime     time.Duration
}

type Question struct {
	Street             geodata.Street
	points             map[string]int
	allPlayersAnswered chan bool
	begin              time.Time
	duration           time.Duration
	number             int
}

func (q *Question) waitForPlayers(countdown func(int)) {
	t2 := time.NewTimer(q.duration - (1 * time.Second))
	t1 := time.NewTimer(q.duration - (2 * time.Second))
	finish := time.NewTimer(q.duration)
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
	if r.StreetList == nil {
		errors = append(errors, "streetListMissing")
	}
	if r.NumberOfQuestions < 1 {
		errors = append(errors, "numberOfQuestionsToSmall")
	}
	if r.NumberOfQuestions > 100 {
		errors = append(errors, "numberOfQuestionsToBig")
	}
	if r.MaxAnswerTime < 10*time.Second {
		errors = append(errors, "maxAnswerTimeToSmall")
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
		creation: time.Now(),
		random:   rand.New(rand.NewSource(time.Now().Unix())),
		players:  make(map[string]*Player),
		options: RoomOptions{
			MaxAnswerTime:     120 * time.Second,
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
		panic(fmt.Sprintf("player with key \"%s\" does not exist in room \"%s\"", playerKey, r.Key()))
	}
	r.options = options
	r.notifyPlayers(
		func(p Player) {
			p.NotifyRoomUpdated(options, player.Name)
		},
	)
}

func (r *roomImpl) Join(name string) Player {
	player := Player{Name: name, Secret: keygen.PlayerKey(), Key: keygen.PlayerKey()}
	r.notifyPlayers(
		func(p Player) {
			p.NotifyPlayerJoined(player.Name, player.Key)
		},
	)
	r.players[player.Key] = &player
	return player
}

func (r *roomImpl) Leave(playerKey string) Player {
	player := r.players[playerKey]
	r.notifyPlayers(
		func(p Player) {
			p.NotifyPlayerLeft(player.Name, player.Key)
		},
	)
	delete(r.players, playerKey)
	if len(r.players) == 0 {
		r.finished = true
	}
	return *player
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
		go consumer(*player)
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
	r.notifyPlayers(
		func(player Player) {
			player.NotifyGameStarted(startPlayer.Key)
		},
	)
	numberOfQuestions := r.options.NumberOfQuestions
	r.points = make(map[string]int)
	go func() {
		r.started = true
		for round := 0; round < numberOfQuestions; round++ {
			err := r.playQuestion(round)
			if err != nil {
				break
			}
			r.advanceGame = make(chan bool)
			<-r.advanceGame
			r.advanceGame = nil
		}
		r.notifyPlayers(
			func(player Player) {
				player.NotifyGameEnded("finished", r.points)
			},
		)
		r.points = nil
		r.finished = true
	}()
}

func (r *roomImpl) playQuestion(round int) error {
	tries := 0
	randomStreet, err := r.options.StreetList.GetRandomStreet(r.random)
	for tries < 10 && err != nil {
		randomStreet, err = r.options.StreetList.GetRandomStreet(r.random)
		tries = tries + 1
	}
	if randomStreet.Coordinate == nil {
		r.notifyPlayers(
			func(player Player) {
				player.NotifyGameEnded("repeatedly failed to get random street", r.points)
			},
		)
		return fmt.Errorf("repeatedly failed to get random street")
	}
	r.Lock()
	r.currentQuestion = &Question{
		Street:             randomStreet,
		points:             make(map[string]int),
		allPlayersAnswered: make(chan bool),
		begin:              time.Now(),
		duration:           r.options.MaxAnswerTime,
		number:             round,
	}
	r.Unlock()
	r.sendCountdowns(
		3, func(player Player) func(int, int) {
			return player.NotifyQuestionCountdown
		}, round,
	)
	r.notifyPlayers(
		func(player Player) {
			player.NotifyQuestion(randomStreet.Name, round)
		},
	)
	r.currentQuestion.waitForPlayers(
		func(followUps int) {
			r.notifyPlayers(
				func(player Player) {
					player.NotifyAnswerTimeCountdown(followUps)
				},
			)
		},
	)
	for key, value := range r.currentQuestion.points {
		r.points[key] = r.points[key] + value
	}
	result := QuestionResult{
		Question:       randomStreet.Name,
		Solution:       *randomStreet.Coordinate,
		PointDelta:     r.currentQuestion.points,
		Points:         r.points,
		QuestionNumber: round,
	}
	r.notifyPlayers(
		func(player Player) {
			player.NotifyQuestionResults(result)
		},
	)
	r.Lock()
	r.currentQuestion = nil
	r.Unlock()
	return nil
}

func (r *roomImpl) AnswerQuestion(playerKey string, guess types.Coordinate) (int, error) {
	_, ok := r.players[playerKey]
	if !ok {
		panic(fmt.Sprintf("player with key \"%s\" not found in this room", playerKey))
	}
	question := r.currentQuestion
	result, err := geodata.VerifyAnswer(guess, question.Street.Name)
	if result {
		difference := time.Now().Sub(question.begin)
		percent := 1.0 * float64(difference.Milliseconds()) / float64(question.duration.Milliseconds())
		question.points[playerKey] = int(math.Max(10, 100-(100*percent)))
	} else {
		question.points[playerKey] = 0
	}
	r.notifyPlayers(
		func(player Player) {
			player.NotifyPlayerAnswered(playerKey, question.points[playerKey])
		},
	)
	if len(question.points) == len(r.players) {
		question.allPlayersAnswered <- true
	}
	return question.points[playerKey], err
}

func (r *roomImpl) HasActiveQuestion(playerKey string) bool {
	if r.currentQuestion == nil {
		return false
	}
	_, ok := r.currentQuestion.points[playerKey]
	return !ok
}

func (r *roomImpl) Question() (string, int) {
	return r.currentQuestion.Street.Name, r.currentQuestion.number
}

func (r *roomImpl) CanBeAdvanced() bool {
	return r.advanceGame != nil
}

func (r *roomImpl) AdvanceToNextQuestion() {
	r.advanceGame <- true
}

func (r *roomImpl) sendCountdowns(
	amount int, consumer func(Player) func(int, int), numberOfQuestion int,
) {
	for i := 0; i < amount; i++ {
		r.notifyPlayers(
			func(player Player) {
				consumer(player)(amount-i-1, numberOfQuestion)
			},
		)
		time.Sleep(time.Second)
	}
}

func (r *roomImpl) Started() bool {
	return r.started
}

func (r *roomImpl) Finished() bool {
	return r.finished
}

func (r *roomImpl) Creation() time.Time {
	return r.creation
}

type Player struct {
	Notifier
	Secret string
	Key    string
	Name   string
}

type QuestionResult struct {
	Question       string           `json:"question"`
	Solution       types.Coordinate `json:"solution"`
	PointDelta     map[string]int   `json:"pointDelta"`
	Points         map[string]int   `json:"points"`
	QuestionNumber int              `json:"questionNumber"`
}

type Notifier interface {
	NotifyPlayerJoined(string, string)
	NotifyPlayerLeft(string, string)
	NotifyRoomUpdated(RoomOptions, string)
	NotifyGameStarted(playerKey string)
	NotifyPlayerAnswered(string, int)
	NotifyQuestionCountdown(int, int)
	NotifyQuestion(string, int)
	NotifyAnswerTimeCountdown(int)
	NotifyQuestionResults(result QuestionResult)
	NotifyGameEnded(reason string, result map[string]int)
}
