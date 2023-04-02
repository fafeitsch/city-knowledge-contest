import { environment } from './environment';
import { catchError, defer, map, Observable, of, switchMap, take } from 'rxjs';
import store from './store';

export function doRpc<ResponseType>(method: string, params: any): Observable<ResponseType> {
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

export async function handleRPCRequest<ResponseType>(method: string, params: any) {
  return fetch(environment[import.meta.env.MODE].apiUrl, {
    method: 'POST',
    body: JSON.stringify({
      method,
      params,
    }),
  })
    .then((response) => response.json())
    .then((response) => {
      if (response.error) {
        throw new Error(response.error.message);
      }
      return response.result;
    });
}

export type RoomConfiguration = {
  listFileName: string;
};

export function updateRoomConfiguration(configuration: RoomConfiguration): Observable<string[]> {
  return store.get.room$.pipe(
    take(1),
    switchMap((authData) =>
      doRpc<{ errors: string[] }>('updateRoom', {
        ...authData,
        ...configuration,
        maxAnswerTimeSec: 30,
        numberOfQuestions: 2,
      }),
    ),
    map((result) => result.errors),
    catchError((err) => {
      console.error(err);
      return of(['noConnection']);
    }),
  );
}

export function answerQuestion(guess: [number, number]): Observable<number> {
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
}
