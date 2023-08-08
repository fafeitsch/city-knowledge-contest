package contest

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/fafeitsch/city-knowledge-contest/backend/geodata"
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"github.com/fafeitsch/city-knowledge-contest/backend/types"
)

type Room struct {
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
	quit            chan bool
}

func (r *Room) Key() string {
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
	quit               chan bool
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
		case <-q.quit:
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

func NewRoom(seedText string) *Room {
	seed := time.Now().Unix()
	if seedText != "" {
		hashFunc := fnv.New32a()
		_, _ = hashFunc.Write([]byte(seedText))
		seed = int64(hashFunc.Sum32())
	}
	return &Room{
		key:      keygen.RoomKey(),
		creation: time.Now(),
		random:   rand.New(rand.NewSource(seed)),
		players:  make(map[string]*Player),
		options: RoomOptions{
			MaxAnswerTime:     120 * time.Second,
			NumberOfQuestions: 10,
		},
		quit: make(chan bool),
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
		panic(fmt.Sprintf("player with key \"%s\" does not exist in room \"%s\"", playerKey, r.Key()))
	}
	r.options = options
	r.notifyPlayers(
		func(p Player) {
			p.NotifyRoomUpdated(options, player.Name)
		},
	)
}

func (r *Room) Join(name string) Player {
	player := Player{Name: name, Secret: keygen.PlayerKey(), Key: keygen.PlayerKey()}
	r.notifyPlayers(
		func(p Player) {
			p.NotifyPlayerJoined(player.Name, player.Key)
		},
	)
	r.players[player.Key] = &player
	return player
}

func (r *Room) Leave(playerKey string) Player {
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

func (r *Room) Kick(target string, initiator string) Player {
	kicked := r.players[target]
	r.notifyPlayers(
		func(p Player) {
			p.NotifyPlayerKicked(kicked.Key, kicked.Name, initiator)
		},
	)
	delete(r.players, target)
	if len(r.players) == 0 {
		r.finished = true
	}
	return *kicked
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
		go consumer(*player)
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
	r.notifyPlayers(
		func(player Player) {
			player.NotifyGameStarted(startPlayer.Key)
		},
	)
	numberOfQuestions := r.options.NumberOfQuestions
	r.points = make(map[string]int)
	go func() {
		r.started = true
	GameLoop:
		for round := 0; round < numberOfQuestions; round++ {
			err := r.playQuestion(round)
			if err != nil {
				break
			}
			r.advanceGame = make(chan bool)
			select {
			case <-r.advanceGame:
				break
			case <-r.quit:
				break GameLoop
			}
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

func (r *Room) Close() error {
	close(r.quit)
	return nil
}

func (r *Room) playQuestion(round int) error {
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
		quit:               r.quit,
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

func (r *Room) AnswerQuestion(playerKey string, guess types.Coordinate) (int, error) {
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

func (r *Room) HasActiveQuestion(playerKey string) bool {
	if r.currentQuestion == nil {
		return false
	}
	_, ok := r.currentQuestion.points[playerKey]
	return !ok
}

func (r *Room) Question() (string, int) {
	return r.currentQuestion.Street.Name, r.currentQuestion.number
}

func (r *Room) CanBeAdvanced() bool {
	return r.advanceGame != nil
}

func (r *Room) AdvanceToNextQuestion() {
	r.advanceGame <- true
}

func (r *Room) sendCountdowns(
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

func (r *Room) Started() bool {
	return r.started
}

func (r *Room) Finished() bool {
	return r.finished
}

func (r *Room) Creation() time.Time {
	return r.creation
}

type Player struct {
	Notifier
	Secret        string
	Key           string
	Name          string
	receivedKicks int
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
	NotifyPlayerKicked(string, string, string)
}
