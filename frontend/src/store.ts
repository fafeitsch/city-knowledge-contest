import { BehaviorSubject, distinctUntilChanged, map } from 'rxjs';

interface State {
  roomId: string | undefined;
  username: string | undefined;
  playerKey: string | undefined;
  gameState: 'inital' | 'waiting' | 'started';
  players: string[];
}

const state: State = {
  roomId: undefined,
  username: undefined,
  playerKey: undefined,
  gameState: 'inital',
  players: [],
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
    gameState$: state$.pipe(
      map((state) => state.gameState),
      distinctUntilChanged()
    ),
    players$: state$.pipe(
      map((state) => state.players),
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
    gameState(gameState: 'inital' | 'waiting' | 'started') {
      state$.next({
        ...state$.value,
        gameState,
      });
    },
    players(players: string[]) {
      state$.next({
        ...state$.value,
        players,
      });
    },
    addPlayer(player: string) {
      const newPlayers = [...state$.value.players, player];
      state$.next({
        ...state$.value,
        players: newPlayers,
      });
    },
  },
};