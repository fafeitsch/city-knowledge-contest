import { BehaviorSubject, distinctUntilChanged, map } from "rxjs";
import Players from "./views/Players.svelte";

export enum GameState {
  SetupUsername = "SetupUsername",
  SetupMap = "SetupMap",
  Waiting = "Waiting",
  Started = "Started",
  QuestionCountdown = "QuestionCountdown",
  Question = "Question",
  Finished = "Finished",
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
  roomId: string | undefined;
  username: string | undefined;
  playerKey: string | undefined;
  playerSecret: string | undefined;
  gameState: GameState;
  players: Player[];
  countdownValue: string;
  question: string;
  game: Game | undefined;
  gameResult: GameResult | undefined;
}

const state: State = {
  roomId: undefined,
  username: undefined,
  playerKey: undefined,
  playerSecret: undefined,
  countdownValue: undefined,
  gameState: GameState.SetupUsername,
  players: [],
  question: undefined,
  game: undefined,
  gameResult: undefined,
};

const state$ = new BehaviorSubject<State>(state);

export default {
  get: {
    roomId$: state$.pipe(
      map((state) => state.roomId),
      distinctUntilChanged()
    ),
    username$: state$.pipe(
      map((state) => state.username),
      distinctUntilChanged()
    ),
    playerKey$: state$.pipe(
      map((state) => state.playerKey),
      distinctUntilChanged()
    ),
    playerSecret$: state$.pipe(
      map((state) => state.playerSecret),
      distinctUntilChanged()
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
  },
  set: {
    roomId(roomId: string | undefined) {
      state$.next({
        ...state$.value,
        roomId,
      });
    },
    username(username: string | undefined) {
      state$.next({
        ...state$.value,
        username,
      });
    },
    playerKey(playerKey: string | undefined) {
      state$.next({
        ...state$.value,
        playerKey,
      });
    },
    playerSecret(playerSecret: string | undefined) {
      state$.next({
        ...state$.value,
        playerSecret,
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
          return { ...player, points: gameResult.points[player.playerKey] };
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
  },
};
