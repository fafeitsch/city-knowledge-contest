import { environment } from "./environment";

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

export async function handleRPCRequest<Params, ResponseType>({
  method,
  params,
}: {
  method: string;
  params: Params;
}) {
  return fetch(environment[import.meta.env.MODE].apiUrl, {
    method: "POST",
    body: JSON.stringify({
      method,
      params,
    }),
  }).then<Response<ResponseType>>((response) => response.json());
}
