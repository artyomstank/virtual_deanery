import { Outlet } from "react-router-dom";
import { Header } from "../../components/layout/Header";
import { Sidebar } from "../../components/layout/Sidebar";

export function AdminLayout() {
  return (
    <div className="min-h-screen flex flex-col bg-zinc-950 text-white">
      <Header />
      <div className="flex flex-1">
        <Sidebar />
        <main className="flex-1 p-8 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
}