import { BehaviorSubject, catchError, distinctUntilChanged, filter, map, of, switchMap, withLatestFrom } from 'rxjs';
import { doRpc, handleRPCRequest } from './rpc';

export type GameState = 'SetupMap' | 'Waiting' | 'Question' | 'Finished' | 'GameEnded';

export type Game = {
  playerKey: string;
  playerSecret: string;
  roomKey: string;
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
  players: Player[];
  countdownValue: number;
  question: string;
  room: Game | undefined;
  gameResult: GameResult | undefined;
  gameConfiguration: GameConfiguration | undefined;
  lastResult: number | undefined;
  gameErrors: string[];
}

export interface GameConfiguration {
  listFileName: string;
  numberOfQuestions: number;
  maxAnswerTimeSec: number;
}

const state: State = {
  countdownValue: undefined,
  players: [],
  question: undefined,
  room: undefined,
  gameResult: undefined,
  lastResult: undefined,
  gameConfiguration: undefined,
  gameErrors: ['notInitialized'],
};

const state$ = new BehaviorSubject<State>(state);

const store = {
  get: {
    gameState$: state$.pipe(
      map<State, GameState>((state) => {
        if (!state.room?.roomKey && !state.gameResult) {
          return 'SetupMap';
        }
        if (
          state.countdownValue === undefined &&
          !state.question &&
          state.lastResult === undefined &&
          !state.gameResult
        ) {
          return 'Waiting';
        }
        if ((state.countdownValue || state.question) && state.lastResult === undefined) {
          return 'Question';
        }
        if (state.lastResult !== undefined) {
          return 'Finished';
        }
        if (state.gameResult && !state.room) {
          return 'GameEnded';
        }
      }),
      distinctUntilChanged(),
    ),
    players$: state$.pipe(
      map((state) => state.players),
      distinctUntilChanged(),
    ),
    countdownValue$: state$.pipe(
      map((state) => state.countdownValue),
      distinctUntilChanged(),
    ),
    question$: state$.pipe(
      map((state) => state.question),
      distinctUntilChanged(),
    ),
    room$: state$.pipe(
      map((state) => state.room),
      distinctUntilChanged(),
    ),
    gameResult$: state$.pipe(
      map((state) => state.gameResult),
      distinctUntilChanged(),
    ),
    lastResult$: state$.pipe(
      map((state) => state.lastResult),
      distinctUntilChanged(),
    ),
    gameConfiguration$: state$.pipe(
      map((state) => state.gameConfiguration),
      distinctUntilChanged(),
    ),
    errors$: state$.pipe(
      map((state) => state.gameErrors),
      distinctUntilChanged(),
    ),
  },
  set: {
    lastResult(result: number) {
      state$.next({
        ...state$.value,
        lastResult: result,
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
        .map((player) => ({
          ...player,
          points: gameResult.points[player.playerKey],
        }))
        .sort((playerA, playerB) => playerA.points - playerB.points);
      state$.next({
        ...state$.value,
        players: newPlayers,
      });
    },
    countdownValue(countdownValue: number) {
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
        room: game,
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
      });
    },
    streetList(fileName: string) {
      state$.next({
        ...state$.value,
        gameConfiguration: {
          ...state$.value.gameConfiguration,
          listFileName: fileName,
        },
      });
    },
    gameErrors(errors: string[]) {
      state$.next({
        ...state$.value,
        gameErrors: errors,
      });
    },
  },
  methods: {
    async startGame() {
      await handleRPCRequest<undefined>('startGame', state$.value.room);
    },
    async createRoom(username: string) {
      const data = await handleRPCRequest<{
        roomKey: string;
        playerKey: string;
        playerSecret: string;
      }>('createRoom', {
        name: username,
      });
      store.set.game(data);
    },
    async joinRoom(roomKey: string, username: string) {
      const data = await handleRPCRequest<{
        playerKey: string;
        playerSecret: string;
      }>('joinRoom', {
        name: username,
        roomKey,
      });
      store.set.game({ ...data, roomKey: roomKey });
    },
    async answerQuestion(guess: number[]) {
      const data = await handleRPCRequest<{ points: number }>('answerQuestion', {
        ...state$.value.room,
        guess,
      });
      store.set.lastResult(data.points);
    },
    async advanceGame() {
      await handleRPCRequest('advanceGame', state$.value.room);
      store.set.lastResult(undefined);
      store.set.question(undefined);
    },
  },
};

export default store;

store.get.gameConfiguration$
  .pipe(
    filter((config) => !!config),
    withLatestFrom(store.get.room$.pipe(filter((authData) => !!authData))),
    switchMap(([config, authData]) =>
      doRpc<{ errors: string[] }>('updateRoom', {
        ...authData,
        ...config,
        maxAnswerTimeSec: 30,
        numberOfQuestions: 2,
      }),
    ),
    map((result) => result.errors),
    catchError((err) => {
      console.error(err);
      return of(['noConnection']);
    }),
  )
  .subscribe(async (result) => {
    store.set.gameErrors(result);
  });
