import { BehaviorSubject, distinctUntilChanged, map, tap } from 'rxjs';

export type Game = {
  playerKey: string;
  playerSecret: string;
  roomKey: string;
};

export type Player = {
  name: string;
  playerKey: string;
  points: number | undefined;
};

interface State {
  players: Player[];
  room: Game | undefined;
}

const state: State = {
  players: [],
  room: undefined,
};

const state$ = new BehaviorSubject<State>(state);

const store = {
  get: {
    players$: state$.pipe(
      map((state) => state.players),
      distinctUntilChanged(),
    ),
    room$: state$.pipe(
      map((state) => state.room),
      distinctUntilChanged(),
    ),
  },
  set: {
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
    updatePlayerRanking(points: Record<string, number | undefined>) {
      const newPlayers = state$.value.players
        .map((player) => ({
          ...player,
          points: points[player.playerKey],
        }))
        .sort((playerA, playerB) => playerA.points - playerB.points);
      state$.next({
        ...state$.value,
        players: newPlayers,
      });
    },
    game(game: Game) {
      state$.next({
        ...state$.value,
        room: game,
      });
    },
    resetGame() {
      state$.next({
        ...state,
      });
    },
  },
};

export default store;
