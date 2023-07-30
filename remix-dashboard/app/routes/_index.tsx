import { json, type V2_MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";

export const meta: V2_MetaFunction = () => [{ title: "Dashboard" }];

export const loader = async () => {
  return json([
    { category: "Data 1", stat: "12,345" },
    { category: "Data 2", stat: "3,456" },
    { category: "Data 3", stat: "67,891" },
  ]);
};

export default function Index() {
  const data = useLoaderData<typeof loader>();
  return (
    <main className="mx-auto max-w-7xl p-10">
      <div className="grid grid-cols-3 gap-8">
        {data.map((d) => (
          <div
            key={d.category}
            className="rounded-lg border-t-4 border-indigo-500 bg-white p-6 shadow ring-1"
          >
            <div className="text-sm font-light text-gray-500">{d.category}</div>
            <div className="text-3xl font-semibold">{d.stat}</div>
          </div>
        ))}
      </div>
    </main>
  );
}
