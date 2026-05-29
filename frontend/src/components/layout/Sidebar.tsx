import {
  BookOpen,
  ClipboardList,
  LayoutDashboard,
  Shield,
  Users,
} from "lucide-react";
import { NavLink } from "react-router-dom";

const items = [
  { label: "Дашборд", path: "/admin", icon: LayoutDashboard },
  { label: "Пользователи", path: "/admin/users", icon: Users },
  { label: "Учебные данные", path: "/admin/data", icon: BookOpen },
  { label: "Отчётность", path: "/admin/reports", icon: ClipboardList },
  { label: "Права доступа", path: "/admin/acl", icon: Shield },
];

export function Sidebar() {
  return (
    <aside className="flex h-[calc(100vh-5rem)] w-16 shrink-0 flex-col border-r border-zinc-800 bg-zinc-900 md:w-64">
      <nav className="flex flex-1 flex-col gap-2 p-3 md:p-4">
        {items.map((item) => {
          const Icon = item.icon;

          return (
            <NavLink
              key={item.path}
              to={item.path}
              end={item.path === "/admin"}
              title={item.label}
              className={({ isActive }) =>
                `flex items-center justify-center gap-3 rounded-lg px-3 py-3 text-sm font-medium transition-all duration-200 md:justify-start md:px-4 ${
                  isActive
                    ? "bg-indigo-600 text-white"
                    : "text-zinc-400 hover:bg-zinc-800 hover:text-white"
                }`
              }
            >
              <Icon size={18} />
              <span className="hidden md:inline">{item.label}</span>
            </NavLink>
          );
        })}
      </nav>
    </aside>
  );
}
