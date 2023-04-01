import { environment } from './environment';
import { defer, map, Observable, switchMap } from 'rxjs';

type Response<T> = {
  jsonrpc: string;
  result: T;
  id: null;
};

export interface UpdateRoomParams {
  listFileName: string;
  numberOfQuestions: number;
  playerKey: string;
  playerSecret: string;
  roomKey: string;
  maxAnswerTimeSec: number;
}

export function doRpc<ResponseType>(
  method: string,
  params: any
): Observable<ResponseType> {
  return defer(() =>
    fetch(environment[import.meta.env.MODE].apiUrl, {
      method: 'POST',
      body: JSON.stringify({
        method,
        params,
      }),
    })
  ).pipe(
    switchMap((response) => response.json()),
    map((response) => {
      if (response.error) {
        throw new Error(response.error.message);
      }
      return response;
    }),
    map((response) => response.result)
  );
}

export async function handleRPCRequest<ResponseType>(
  method: string,
  params: any
) {
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
