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
  delta: number | undefined;
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
    removePlayer(removedPlayer: Player) {
      const newPlayers: Player[] = state$.value.players.filter(player => player.playerKey !== removedPlayer.playerKey);
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
          delta: undefined,
        }))
        .sort((playerA, playerB) => playerB.points - playerA.points);
      state$.next({
        ...state$.value,
        players: newPlayers,
      });
    },
    updatePlayerDelta(payload: { playerKey: string; pointsDelta: number }) {
      const newPlayers = state$.value.players
        .map((player) => {
          if (payload.playerKey === player.playerKey) {
            return {
              ...player,
              delta: payload.pointsDelta,
            };
          }
          return player;
        })
        .sort((playerA, playerB) => playerB.points + playerB.delta - (playerA.points + playerA.delta));
      state$.next({
        ...state$.value,
        players: newPlayers,
      });
    },
    removePlayerDelta() {
      state$.next({
        ...state$.value,
        players: state$.value.players.map((player) => ({ ...player, delta: undefined })),
      });
    },
    game(game: Game) {
      state$.next({
        ...state$.value,
        room: { roomKey: game.roomKey, playerKey: game.playerKey, playerSecret: game.playerSecret },
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
