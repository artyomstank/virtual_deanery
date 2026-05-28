import {
  LayoutDashboard,
  Users,
  GraduationCap,
  CalendarDays,
  Shield,
  ClipboardList,
  BookOpen,
  UserSquare2
} from "lucide-react";

import { NavLink } from "react-router-dom";

const items = [
  {
    label: "Dashboard",
    path: "/dashboard",
    icon: LayoutDashboard
  },
  {
    label: "Users",
    path: "/users",
    icon: Users
  },
  {
    label: "Students",
    path: "/students",
    icon: GraduationCap
  },
  {
    label: "Teachers",
    path: "/teachers",
    icon: UserSquare2
  },
  {
    label: "Schedule",
    path: "/schedule",
    icon: CalendarDays
  },
  {
    label: "Grades",
    path: "/grades",
    icon: BookOpen
  },
  {
    label: "ACL",
    path: "/acl",
    icon: Shield
  },
  {
    label: "Audit",
    path: "/audit",
    icon: ClipboardList
  }
];

export function Sidebar() {
  return (
    <aside className="hidden w-72 border-r border-zinc-800 bg-zinc-900/50 lg:flex lg:flex-col">
      <div className="border-b border-zinc-800 px-6 py-6">
        <h1 className="text-xl font-bold">
          University Admin
        </h1>

        <p className="mt-1 text-sm text-zinc-400">
          Internal management system
        </p>
      </div>

      <nav className="flex flex-1 flex-col gap-2 p-4">
        {items.map((item) => {
          const Icon = item.icon;

          return (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `
                flex items-center gap-3 rounded-xl px-4 py-3
                transition-all duration-200

                ${
                  isActive
                    ? "bg-indigo-500 text-white"
                    : "text-zinc-400 hover:bg-zinc-800 hover:text-white"
                }
              `
              }
            >
              <Icon size={18} />

              <span className="font-medium">
                {item.label}
              </span>
            </NavLink>
          );
        })}
      </nav>
    </aside>
  );
}