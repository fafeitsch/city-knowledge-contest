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
  return fetch("http://localhost:23123/rpc", {
    method: "POST",
    body: JSON.stringify({
      method,
      params,
    }),
  }).then<Response<ResponseType>>((response) => response.json());
}
