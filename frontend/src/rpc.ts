import { environment } from "./environment";

type Response<T> = {
  jsonrpc: string;
  result: T;
  id: null;
};

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
