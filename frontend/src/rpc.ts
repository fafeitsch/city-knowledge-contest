import { environment } from './environment';
import { catchError, defer, EMPTY, map, Observable, of, switchMap, take, tap } from 'rxjs';
import store from './store';

function doRpc<ResponseType>(method: string, params: any): Observable<ResponseType> {
  return defer(() =>
    fetch(environment[import.meta.env.MODE].apiUrl, {
      method: 'POST',
      body: JSON.stringify({
        method,
        params,
      }),
    }),
  ).pipe(
    switchMap((response) => response.json()),
    map((response) => {
      if (response.error) {
        throw new Error(response.error.message);
      }
      return response;
    }),
    map((response) => response.result),
  );
}

const rpc = {
  updateRoomConfiguration(configuration: RoomConfiguration): Observable<string[]> {
    return store.get.room$.pipe(
      take(1),
      switchMap((authData) =>
        doRpc<{ errors: string[] }>('updateRoom', {
          ...authData,
          ...configuration,
        }),
      ),
      map((result) => result.errors),
      catchError((err) => {
        console.error(err);
        return of(['noConnection']);
      }),
    );
  },
  answerQuestion(guess: [number, number]): Observable<number> {
    return store.get.room$.pipe(
      take(1),
      switchMap((authData) =>
        doRpc<{ points: number }>('answerQuestion', {
          ...authData,
          guess,
        }),
      ),
      map((result) => result.points),
      catchError((err) => {
        console.error(err);
        return of(0);
      }),
    );
  },
  advanceRoom(): Observable<void> {
    return store.get.room$.pipe(
      take(1),
      switchMap((authData) =>
        doRpc<void>('advanceGame', {
          ...authData,
        }),
      ),
      catchError((err) => {
        console.error(err);
        return EMPTY;
      }),
    );
  },
  startGame(): Observable<void> {
    return store.get.room$.pipe(
      take(1),
      switchMap((authData) =>
        doRpc<void>('startGame', {
          ...authData,
        }),
      ),
      catchError((err) => {
        console.error(err);
        return EMPTY;
      }),
    );
  },
  createRoom(username: string): Observable<Room> {
    return doRpc<Room>('createRoom', {
      name: username,
    });
  },
  joinRoom(roomKey: string, username: string): Observable<Room> {
    return doRpc<{
      playerKey: string;
      playerSecret: string;
    }>('joinRoom', {
      roomKey,
      name: username,
    }).pipe(map((data) => ({ ...data, roomKey })));
  },
  getStreetLists(): Observable<StreetList[]> {
    return doRpc<StreetList[]>('getAvailableStreetLists', {});
  },
  getLegalInformation(): Observable<{
    imprint: string;
    dataProtection: string;
    tileServer: string;
    nominatimServer: string;
  }> {
    return doRpc('getLegalInformation', {});
  },
};

export default rpc;

export type RoomConfiguration = {
  listFileName: string;
  numberOfQuestions: number;
  maxAnswerTimeSec: number;
};

export type Room = {
  playerKey: string;
  playerSecret: string;
  roomKey: string;
};

type StreetList = {
  FileName: string;
  name: string;
  center: { Lat: number; Lng: number };
};
