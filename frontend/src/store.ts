import { BehaviorSubject, distinctUntilChanged, map } from 'rxjs';

interface State {
  roomId: string | undefined;
  username: string | undefined;
  gameState: 'inital' | 'waiting' | 'started';
}

const state: State = {
  roomId: undefined,
  username: undefined,
  gameState: 'inital',
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
    gameState$: state$.pipe(
      map((state) => state.gameState),
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
    gameState(gameState: 'inital' | 'waiting' | 'started') {
      state$.next({
        ...state$.value,
        gameState,
      });
    },
  },
};
