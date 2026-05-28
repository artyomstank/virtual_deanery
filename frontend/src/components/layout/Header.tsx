import { Bell, Search } from "lucide-react";

export function Header() {
  return (
    <header className="flex h-20 items-center justify-between border-b border-zinc-800 bg-zinc-900/50 px-6 backdrop-blur-xl">
      <div className="flex items-center gap-3 rounded-xl border border-zinc-800 bg-zinc-900 px-4 py-3 text-sm text-zinc-400">
        <Search size={16} />

        <input
          placeholder="Search..."
          className="bg-transparent outline-none"
        />
      </div>

      <div className="flex items-center gap-4">
        <button className="flex h-11 w-11 items-center justify-center rounded-xl border border-zinc-800 bg-zinc-900 hover:bg-zinc-800">
          <Bell size={18} />
        </button>

        <div className="flex items-center gap-3 rounded-xl border border-zinc-800 bg-zinc-900 px-3 py-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-indigo-500 font-semibold">
            A
          </div>

          <div>
            <p className="text-sm font-medium">
              Admin
            </p>

            <p className="text-xs text-zinc-400">
              administrator
            </p>
          </div>
        </div>
      </div>
    </header>
  );
}