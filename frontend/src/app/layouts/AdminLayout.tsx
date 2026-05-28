import { Outlet } from "react-router-dom";

export function AdminLayout() {
  return (
    <div className="flex min-h-screen bg-zinc-950 text-white">
      <aside className="w-64 border-r border-zinc-800">
        Sidebar placeholder
      </aside>

      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  );
}