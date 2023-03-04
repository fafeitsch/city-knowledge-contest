import {BehaviorSubject, distinctUntilChanged, filter, map, withLatestFrom} from "rxjs";
import {handleRPCRequest, type UpdateRoomParams} from './rpc';

export enum GameState {
  SetupUsername = "SetupUsername",
  SetupMap = "SetupMap",
  Waiting = "Waiting",
  Started = "Started",
  QuestionCountdown = "QuestionCountdown",
  Question = "Question",
  Finished = "Finished",
  GameEnded = "GameEnded",
}

export type Game = {
  playerKey: string;
  playerSecret: string;
  roomId: string;
};

export type Player = {
  name: string;
  playerKey: string;
  points: number | undefined;
  delta: number | undefined;
};

export type GameResult = {
  delta: Record<string, number>;
  points: Record<string, number>;
  question: string;
  solution: {
    lat: number;
    lon: number;
  };
};

interface State {
  username: string;
  gameState: GameState;
  players: Player[];
  countdownValue: string;
  question: string;
  game: Game | undefined;
  gameResult: GameResult | undefined;
  gameConfiguration: GameConfiguration | undefined;
  gameErrors: string[]
}

export interface GameConfiguration {
  listFileName: string;
  numberOfQuestions: number;
  maxAnswerTimeSec: number;
}

const state: State = {
  username: undefined,
  countdownValue: undefined,
  gameState: GameState.SetupUsername,
  players: [],
  question: undefined,
  game: undefined,
  gameResult: undefined,
  gameConfiguration: undefined,
  gameErrors: ['notInitialized'],
};

const state$ = new BehaviorSubject<State>(state);

const store = {
  get: {
    authData$: state$.pipe(
      map(state => state.game ? ({
        playerKey: state.game.playerKey,
        playerSecret: state.game.playerSecret,
        roomKey: state.game.roomId
      }) : undefined),
      distinctUntilChanged(),
    ),
    gameState$: state$.pipe(
      map((state) => state.gameState),
      distinctUntilChanged()
    ),
    players$: state$.pipe(
      map((state) => state.players),
      distinctUntilChanged()
    ),
    countdownValue$: state$.pipe(
      map((state) => state.countdownValue),
      distinctUntilChanged()
    ),
    question$: state$.pipe(
      map((state) => state.question),
      distinctUntilChanged()
    ),
    game$: state$.pipe(
      map((state) => state.game),
      distinctUntilChanged()
    ),
    gameResult$: state$.pipe(
      map((state) => state.gameResult),
      distinctUntilChanged()
    ),
    gameConfiguration$: state$.pipe(
      map((state) => state.gameConfiguration),
      distinctUntilChanged()
    ),
    errors$: state$.pipe(
      map(state =>state.gameErrors),
      distinctUntilChanged()
    )
  },
  set: {
    username(username: string | undefined) {
      state$.next({
        ...state$.value,
        username,
      });
    },
    gameState(gameState: GameState) {
      state$.next({
        ...state$.value,
        gameState,
      });
    },
    players(players: Player[]) {
      state$.next({
        ...state$.value,
        players,
      });
    },
    addPlayer(player: Player) {
      const newPlayers = [...state$.value.players, player];
      state$.next({
        ...state$.value,
        players: newPlayers,
      });
    },
    updatePlayerRanking(gameResult: GameResult) {
      const newPlayers = state$.value.players
        .map((player) => {
          return {...player, points: gameResult.points[player.playerKey]};
        })
        .sort((playerA, playerB) => {
          return playerA.points > playerB.points ? -1 : 1;
        });
      state$.next({
        ...state$.value,
        players: newPlayers,
      });
    },
    countdownValue(countdownValue: string) {
      state$.next({
        ...state$.value,
        countdownValue,
      });
    },
    question(question: string) {
      state$.next({
        ...state$.value,
        question,
      });
    },
    game(game: Game) {
      state$.next({
        ...state$.value,
        game,
      });
    },
    gameResult(gameResult: GameResult) {
      state$.next({
        ...state$.value,
        gameResult,
      });
    },
    resetGame() {
      state$.next({
        ...state,
        username: state$.value.username,
        gameState: GameState.SetupMap,
      });
    },
    streetList(fileName: string) {
      state$.next({
        ...state$.value,
        gameConfiguration: {
          ...state$.value.gameConfiguration,
          listFileName: fileName,
        }
      })
    },
    gameErrors(errors: string[]) {
      state$.next({
       ...state$.value,
        gameErrors: errors,
      });
    }
  },
  methods: {
    async startGame() {
      await handleRPCRequest<
        { playerKey: string; playerSecret: string; roomKey: string },
        null
      >({
        method: "startGame",
        params: {
          playerKey: state$.value.game.playerKey,
          roomKey: state$.value.game.roomId,
          playerSecret: state$.value.game.playerSecret,
        },
      })
      store.set.gameState(GameState.Started);
    }
  }
};

export default store

store.get.gameConfiguration$.pipe(
  withLatestFrom(store.get.authData$.pipe(filter(authData => !!(authData)))),
).subscribe(
  async ([config, authData]) => {
    const result = await handleRPCRequest<UpdateRoomParams, {errors: string[]}>({
      method: "updateRoom",
      params: {
        ...authData,
        ...config,
        maxAnswerTimeSec: 30,
        numberOfQuestions: 10
      },
    })
    store.set.gameErrors(result.result.errors);
  }
)
