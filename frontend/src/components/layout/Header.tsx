import { LogOut } from "lucide-react";
import { useAuthStore } from "../../store/auth.store";
import { Button } from "../ui/Form";

const roleLabels: Record<string, string> = {
  student: "студент",
  teacher: "преподаватель",
  dean: "декан",
  admin: "админ",
};

export function Header() {
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);

  const handleLogout = () => {
    logout();
    window.location.href = "/login";
  };

  return (
    <header className="border-b border-zinc-800 bg-zinc-950 px-4 py-4 md:px-6">
      <div className="flex items-center justify-between gap-4">
        <h1 className="text-xl font-bold text-white md:text-2xl">
          Виртуальный деканат
        </h1>
        <div className="flex items-center gap-3">
          <span className="hidden text-sm text-zinc-400 sm:inline">
            {user?.name || "Администратор"}
            {user?.role ? ` (${roleLabels[user.role] || user.role})` : ""}
          </span>
          <Button
            variant="ghost"
            size="sm"
            onClick={handleLogout}
            className="flex items-center gap-2"
          >
            <LogOut size={16} />
            <span className="hidden sm:inline">Выйти</span>
          </Button>
        </div>
      </div>
    </header>
  );
}
